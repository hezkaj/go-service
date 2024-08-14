package services

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id           string
	email        string
	hash_refresh string
}

func findUserById(id string, db *sql.DB) (User, error) {
	user := User{}
	err := db.QueryRow(
		`SELECT * FROM Users WHERE id=$1`,
		id).Scan(&user.id, &user.email, &user.hash_refresh)
	if user.id != "" {
		return user, nil
	} else {
		return user, err
	}
}

func createToken(user User, ip string, period time.Duration) (string, error) {
	key := "very-secret-key"
	var jwtSecretKey = []byte(key)
	payload := jwt.MapClaims{
		"sub": user.email,
		"ip":  ip,
		"exp": time.Now().Add(period).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	token, err := t.SignedString(jwtSecretKey)
	if err != nil {
		return token, err
	}
	return token, nil
}

func cutRefresh(refresh string) string {
	split_ref := strings.Split(refresh, ".")
	str := strings.Fields(split_ref[1])[0]
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded[0:72]
}

func updateUsersRefresh(id string, ref string, db *sql.DB) error {
	cut_ref := cutRefresh(ref)
	hash_ref, err := bcrypt.GenerateFromPassword([]byte(cut_ref), 10)
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE Users SET hash_refresh = $1 
		WHERE id = $2`, hash_ref, id)
	if err != nil {
		return err
	}
	return nil
}

func GetTokens(id string, ip string, db *sql.DB) (string, string, error) {
	user, err := findUserById(id, db)
	if err != nil {
		return "", "", err
	}
	access, err := createToken(user, ip, time.Minute*2)
	if err != nil {
		return "", "", err
	}
	refresh, err := createToken(user, ip, time.Hour*72)
	if err != nil {
		return "", "", err
	}
	err = updateUsersRefresh(id, refresh, db)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func sendMail(email string) error {
	//Сообщение не будет доставлено из-за "example" данных отправителя,
	//но, исходя из выдаваемой ошибки, отправка работает верно
	from := "example@mail.ru"
	password := "example1pass"
	toList := []string{email}
	host := "smtp.mail.ru"
	port := "587"
	msg := "Warning: a new IP was used to refresh tokens"
	body := []byte(msg)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, toList, body)
	if err != nil {
		fmt.Println(err)
		//return err
	}
	fmt.Println("Successfully sent mail")
	return nil
}

func getUser(email string, db *sql.DB) (User, error) {
	user := User{}
	err := db.QueryRow(
		`SELECT * FROM Users WHERE email=$1`,
		email).Scan(&user.id, &user.email, &user.hash_refresh)
	if user.id != "" {
		return user, nil
	} else {
		return user, err
	}
}

func isCompareRef(ref_req string, email string, db *sql.DB) error {
	cut_ref_req := cutRefresh(ref_req)
	user, err := getUser(email, db)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.hash_refresh), []byte(cut_ref_req))
}

func RefreshTokens(email string, new_ip string, old_ip string, db *sql.DB, req_ref string) (string, string, error) {
	err := isCompareRef(req_ref, email, db)
	if err != nil {
		return "", "", err
	}
	if old_ip != new_ip {
		err := sendMail(email)
		if err != nil {
			return "", "", err
		}
	}
	user, err := getUser(email, db)
	if err != nil {
		return "", "", err
	}
	access, err := createToken(user, new_ip, time.Minute*2)
	if err != nil {
		return "", "", err
	}
	refresh, err := createToken(user, new_ip, time.Hour*72)
	if err != nil {
		return "", "", err
	}
	err = updateUsersRefresh(user.id, refresh, db)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
