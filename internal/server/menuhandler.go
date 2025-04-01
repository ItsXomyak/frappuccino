package server

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
	"io"
	"net/http"
	"strconv"

	converter "frappuccino/pkg/convertor"
)

func (h *Handler) CreateNewMenuItem(w http.ResponseWriter, r *http.Request) {
	statusCode := 201
	text := "menu created"

	defer func() {
		Respond(w, statusCode, text)
	}()

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		statusCode = 400
		text = "Invalid request body"
		return
	}

	var newItem models.MenuItem
	err = json.Unmarshal(data, &newItem)
	if err != nil {
		statusCode = 400
		text = "Invalid JSON format"
		return
	}

	_, err = h.Service.CreateMenuItem(newItem)
	if err != nil {
		if errors.Is(err, cerrors.ErrExist) {
			statusCode = 409
			text = "Menu item already exists"
			return
		}
		statusCode = 400
		text = err.Error()
		return
	}
}

func (h *Handler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "all menu got"

	defer func() {
		Respond(w, statusCode, text)
	}()

	data, err := h.Service.GetAllMenuItems()
	if err != nil {
		statusCode = 400
		text = err.Error()
		return
	}

	responseData, err := json.Marshal(data)
	if err != nil {
		statusCode = 400
		text = err.Error()
		return
	}

	w.Write(responseData)
}

func (h *Handler) GetMenuItemByID(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "menu got"

	defer func() {
		Respond(w, statusCode, text)
	}()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode = 400
		text = "Invalid order ID: must be an integer"
		return
	}

	data, err := h.Service.GetMenuItemByID(id)
	if err != nil {
		if errors.Is(err, cerrors.ErrMenuItemNotFound) {
			statusCode = 404
			text = cerrors.ErrMenuItemNotFound.Error()
			return
		}
		statusCode = 400
		text = err.Error()
		return
	}

	responseData, err := json.Marshal(data)
	if err != nil {
		statusCode = 400
		text = err.Error()
		return
	}

	w.Write(responseData)
}

func (h *Handler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "menu updated"

	defer func() {
		Respond(w, statusCode, text)
	}()

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		statusCode = 400
		text = err.Error()
		return
	}

	var updatedItem models.MenuItem
	if err := json.Unmarshal(data, &updatedItem); err != nil {
		statusCode = 400
		text = err.Error()
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode = 400
		text = "Invalid order ID: must be an integer"
		return
	}

	_, err = h.Service.UpdateMenuItem(id, updatedItem)
	if err != nil {
		switch {
		case errors.Is(err, cerrors.ErrMenuItemNotFound):
			statusCode = 404
			text = cerrors.ErrMenuItemNotFound.Error()
		default:
			statusCode = 400
			text = err.Error()
		}
		return
	}
}

func (h *Handler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	statusCode := 204
	text := "menu deleted"

	defer func() {
		Respond(w, statusCode, text)
	}()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode = 400
		text = "Invalid order ID: must be an integer"
		return
	}

	if err := h.Service.DeleteMenuItem(id); err != nil {
		text = err.Error()

		if errors.Is(err, cerrors.ErrMenuItemNotFound) {
			statusCode = 404
			return
		}
		statusCode = 400
		return
	}
}

func Respond(w http.ResponseWriter, statusCode int, text string) {
	w.WriteHeader(statusCode)
	str := converter.Wrap(statusCode, text)
	bb, err := json.Marshal(str)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(statusCode)
	w.Write(bb)
}
