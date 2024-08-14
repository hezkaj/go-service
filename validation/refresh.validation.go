package validation

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func CheckValidAndGetMailAndIP(tokenString string) (error, string, string) {
	var ip string
	var email string
	var mySigningKey = []byte("very-secret-key")
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неожиданный метод подписи: %v", t.Header["alg"])
		}
		return mySigningKey, nil
	}
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	for key, val := range claims {
		if key == "ip" {
			ip = fmt.Sprintf("%v", val)
		} else if key == "sub" {
			email = fmt.Sprintf("%v", val)
		}
	}
	if err != nil {
		return err, email, ip
	}
	if !parsedToken.Valid {
		return err, email, ip
	}
	return nil, email, ip
}
