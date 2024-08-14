package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	id           string
	email        string
	hash_refresh string
}

func DatabaseConnection() (*sql.DB, error) {
	//Пдключемся к базе данных:
	connStr := "user=postgres password=root dbname=godb sslmode=disable host=postgres"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to the database")
	//Создаем таблицу если она не существует:
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			hash_refresh VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}
	//Создаем пользователя, если он не найден в таблице:
	p := User{}
	db.QueryRow(`SELECT * FROM Users WHERE email=$1`, "someone@mail.ru").Scan(&p.id, &p.email, &p.hash_refresh)
	if p.id == "" {
		_, err = db.Exec(`
			INSERT INTO users (id, email, hash_refresh)
			VALUES ('123e4567-e89b-12d3-a456-426614174000', 'someone@mail.ru', 'notusedonlyforfirstrequest');
		`)
		if err != nil {
			return nil, err
		}
		fmt.Println("User inserted successfully")
	}
	return db, nil
}
