package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing/entity"
	"testing/service"

	"gorm.io/gorm"
)

type AuthorAPI struct {
	authorService *service.AuthorService
}

func NewAuthorAPI(
	authorService *service.AuthorService,
) *AuthorAPI {
	return &AuthorAPI{
		authorService: authorService,
	}
}

func (a *AuthorAPI) GetAllAuthor(w http.ResponseWriter, r *http.Request) {

	author := r.URL.Query()
	authorID, foundAuthorId := author["author_id"]

	if foundAuthorId {
		aID, _ := strconv.Atoi(authorID[0])
		authorByID, err := a.authorService.GetAuthorByID(r.Context(), aID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if authorByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error author not found"))
			return
		}

		WriteJSON(w, http.StatusOK, authorByID)
		return
	}

	list, err := a.authorService.GetAllAuthor(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (a *AuthorAPI) CreateNewAuthor(w http.ResponseWriter, r *http.Request) {
	var author entity.AuthorReg

	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid author request"))
		return
	}

	if author.Name == "" || author.Country == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid author's data request"))
		return
	}

	_, err = a.authorService.AddAuthor(r.Context(), entity.Author{
		Name:    author.Name,
		Country: author.Country,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"message": "success create new author",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (a *AuthorAPI) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	aID, _ := strconv.Atoi(authorID)
	err := a.authorService.DeleteAuthor(r.Context(), aID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"author_id": aID,
		"message":   "success delete author's data",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (a *AuthorAPI) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	var author entity.AuthorReg

	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid author request"))
		return
	}

	id := r.URL.Query().Get("author_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid author id"))
		return
	}

	authors, err := a.authorService.UpdateAuthor(r.Context(), entity.Author{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		Name:    author.Name,
		Country: author.Country,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"author_id": authors.ID,
		"message":   "success update author's data",
		"data":      authors,
	}

	WriteJSON(w, http.StatusOK, response)
}
