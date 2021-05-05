package book

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/supachai/api/bookapi/cors"
)

const booksPath = "books"

func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		BooksList, err := getAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(BooksList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var book Books
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = addBook(book)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleBook(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", booksPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bookID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		book, err := getOneBook(bookID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if book == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPut:
		var book Books
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if book.Id != bookID { // from pluralsight have anterisk at book.Id (*book.Id != bookID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateBook(book)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		err := deleteBook(bookID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SetupRoutes(apiBasePath string) {
	booksHandler := http.HandlerFunc(handleBooks)
	bookHandler := http.HandlerFunc(handleBook)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, booksPath), cors.Middleware(booksHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, booksPath), cors.Middleware(bookHandler))
}
