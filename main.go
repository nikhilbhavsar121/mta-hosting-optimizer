package main

import (
	"fmt"
	"log"
	"mta-hosting-optimizer/routes"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

func main() {
	LoadConfig()

	routes.SetupRoutes()
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
}
