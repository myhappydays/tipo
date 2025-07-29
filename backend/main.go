package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"tipo-backend/handlers"
)

func main() {
	// .env 파일 로드
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/api/content", handlers.ContentHandler)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
