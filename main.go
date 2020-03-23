package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type book struct {
	Id   string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type bookArray []book

var books = bookArray{
	{
		Id:   "1",
		Title:  "Golang Cookbook",
		Author: "John Doe",
	},
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook book
	// Convert r.Body into a readable format
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error: Create Book failed")
	}

	json.Unmarshal(reqBody, &newBook)

	// Add the newly created book to the array of books
	books = append(books, newBook)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created book
	json.NewEncoder(w).Encode(newBook)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	bookId := mux.Vars(r)["id"]

	// Get the details from an existing book
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleBook := range books {
		if singleBook.Id == bookId {
			json.NewEncoder(w).Encode(singleBook)
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	bookId := mux.Vars(r)["id"]
	var updatedBook book
	// Convert r.Body into a readable format
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error: Update Book with id=%v failed", bookId)
	}

	json.Unmarshal(reqBody, &updatedBook)

	for i, singleBook := range books {
		if singleBook.Id == bookId {
			singleBook.Title = updatedBook.Title
			singleBook.Author = updatedBook.Author
			books[i] = singleBook
			json.NewEncoder(w).Encode(singleBook)
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	bookId := mux.Vars(r)["id"]

	for i, singleBook := range books {
		if singleBook.Id == bookId {
			// delete book
			books = append(books[:i], books[i+1:]...)
			w.WriteHeader(http.StatusAccepted)
		}
	}
}