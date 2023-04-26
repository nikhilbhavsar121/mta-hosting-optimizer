package main

import (
	"fmt"
	"log"
	"mta-hosting-optimizer/routes"
	"mta-hosting-optimizer/utils"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

func main() {
	LoadConfig()

	conn, dbErr := utils.GetDBConnection()

	if dbErr != nil {
		fmt.Println("Cannot connect to database", dbErr)
		return
	}

	defer conn.Close()

	routes.SetupRoutes(conn)
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
}
