package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing/entity"
	"testing/service"
	"time"

	"gorm.io/gorm"
)

type BookAPI struct {
	bookService *service.BookService
}

func NewBookAPI(
	bookService *service.BookService,
) *BookAPI {
	return &BookAPI{
		bookService: bookService,
	}
}

func (b *BookAPI) GetAllBook(w http.ResponseWriter, r *http.Request) {

	book := r.URL.Query()
	bookID, foundBookId := book["book_id"]

	if foundBookId {
		bID, _ := strconv.Atoi(bookID[0])
		bookByID, err := b.bookService.GetBookByID(r.Context(), bID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if bookByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error book not found"))
			return
		}

		WriteJSON(w, http.StatusOK, bookByID)
		return
	}

	list, err := b.bookService.GetAllBook(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (b *BookAPI) CreateNewBook(w http.ResponseWriter, r *http.Request) {
	var book entity.BookReq

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid book request"))
		return
	}

	if book.Title == "" || book.ISBN == "" || book.PublishedYear == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid book's data request"))
		return
	}

	_, err = b.bookService.AddBook(r.Context(), book)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"message": "success create new book",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (b *BookAPI) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.URL.Query().Get("book_id")
	bID, _ := strconv.Atoi(bookID)
	err := b.bookService.DeleteBook(r.Context(), bID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"book_id": bID,
		"message": "success delete book's data",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (b *BookAPI) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book entity.BookReq

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid book request"))
		return
	}

	t, _ := time.Parse("2006-01-02", book.PublishedYear)

	id := r.URL.Query().Get("book_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid book id"))
		return
	}

	books, err := b.bookService.UpdateBook(r.Context(), entity.Book{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		Title:         book.Title,
		PublishedYear: t,
		ISBN:          book.ISBN,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"book_id": books.ID,
		"message": "success update book's data",
		"data":    books,
	}

	WriteJSON(w, http.StatusOK, response)
}
