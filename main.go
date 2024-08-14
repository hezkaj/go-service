package main

import (
	"ai-test/controller"
	"ai-test/database"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := database.DatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database initialized successfully")

	controller.Router(db)

	fmt.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
