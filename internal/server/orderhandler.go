package server

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) AddNewOrder(w http.ResponseWriter, r *http.Request) {
	statusCode := 201
	text := "Order created successfully"

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

	var newItems models.Order
	err = json.Unmarshal(data, &newItems)
	if err != nil {
		statusCode = 400
		text = "Invalid JSON format"
		return
	}

	if err := h.Service.OrderCreate(newItems); err != nil {
		if errors.Is(err, cerrors.ErrExist) {
			statusCode = 409
			text = cerrors.ErrExist.Error()
			return
		}
		statusCode = 400
		text = err.Error()
		return
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	statusCode := 201
	text := "Order got successfully"

	defer func() {
		Respond(w, statusCode, text)
	}()

	data, err := h.Service.Get()
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

func (h *Handler) GetOrderId(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "Order id got successfully"

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

	data, err := h.Service.GetId(id)
	if err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = cerrors.ErrNotExist.Error()
		} else {
			statusCode = 400
			text = err.Error()
		}
		return
	}

	responseData, err := json.Marshal(data)
	if err != nil {
		statusCode = 400
		text = "Error during data serialization"
		return
	}

	w.Write(responseData)
}

func (h *Handler) UpDateOrder(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "order updated"

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

	var newItems models.Order
	err = json.Unmarshal(data, &newItems)
	if err != nil {
		statusCode = 400
		text = "Invalid JSON format"
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode = 400
		text = "Invalid order ID: must be an integer"
		return
	}

	if err := h.Service.Update(id, newItems); err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = cerrors.ErrNotExist.Error()
		}
		statusCode = 400
		text = err.Error()
		return
	}
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	statusCode := 204
	text := "order deleted"

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

	if err := h.Service.RemoveOrder(id); err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = cerrors.ErrNotExist.Error()
		} else {
			statusCode = 400
			text = err.Error()
		}
		return
	}
}

func (h *Handler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "order closed"

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

	if err := h.Service.CloseOrder(id); err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = err.Error()
		}

		statusCode = 500
		text = err.Error()
		return
	}
}

func (h *Handler) GetNumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	results, err := h.Service.GetNumberOfOrderedItems(startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to get number of ordered items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) BatchProcessOrders(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Orders []models.Order `json:"orders"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.Service.BatchProcessOrders(request.Orders)
	if err != nil {
		http.Error(w, "Failed to process orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
