package book

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/supachai/api/bookapi/database"
)

// retrieve all books
func getAllBooks() ([]Books, error) {
	results, err := database.DbConn.Query(`SELECT * FROM books`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	books := make([]Books, 0)
	for results.Next() {
		var book Books
		results.Scan(&book.Id, &book.Name, &book.Author, &book.Created_at, &book.Updated_at)
		books = append(books, book)
	}
	return books, nil
}

func getOneBook(bookID int) (*Books, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results := database.DbConn.QueryRowContext(ctx, `
	SELECT * 
	FROM books
	WHERE Id = ?`, bookID)

	book := &Books{}
	err := results.Scan(
		&book.Id,
		&book.Name,
		&book.Author,
		&book.Updated_at,
		&book.Created_at,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return book, nil
}

// add a new book
func addBook(book Books) error {
	_, err := database.DbConn.Query(`INSERT INTO books (name, author) 
		VALUES (?, ?)`,
		book.Name,
		book.Author)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func updateBook(book Books) error {
	if book.Id == 0 { // || *book.Id == 0
		return errors.New("book has invalid ID")
	}
	_, err := database.DbConn.Exec(`UPDATE books SET
		name =?,
		author =?`,
		book.Name, book.Author)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func deleteBook(bookID int) error {
	_, err := database.DbConn.Exec(`DELETE FROM books WHERE Id = ?`, bookID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
