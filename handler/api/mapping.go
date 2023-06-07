package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing/entity"
	"testing/service"
)

type MappingAPI struct {
	mappingService *service.MappingService
}

func NewMappingAPI(
	mappingService *service.MappingService,
) *MappingAPI {
	return &MappingAPI{
		mappingService: mappingService,
	}
}

func (a *MappingAPI) GetAllMapping(w http.ResponseWriter, r *http.Request) {
	list, err := a.mappingService.GetAllMapping(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (a *MappingAPI) CreateNewMapping(w http.ResponseWriter, r *http.Request) {
	var mapping entity.Mapping

	err := json.NewDecoder(r.Body).Decode(&mapping)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid mapping request"))
		return
	}
	_, err = a.mappingService.AddMapping(r.Context(), mapping)
	if err != nil {
		if errors.Is(err, service.ErrDataMapping) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		}
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"message": "success create new mapping",
	}

	WriteJSON(w, http.StatusCreated, response)
}
