package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/supachai/api/bookapi/book"
	"github.com/supachai/api/bookapi/database"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	book.SetupRoutes(basePath)

	// router.HandleFunc("/books", handleBooks)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
