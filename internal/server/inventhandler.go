package server

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/models"
	"frappuccino/internal/svc"
	"frappuccino/pkg/cerrors"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	Service svc.Svc
}

func New(s svc.Svc) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "inventory got"

	defer func() {
		Respond(w, statusCode, text)
	}()

	data, err := h.Service.InventoriesGet()
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

func (h *Handler) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	statusCode := 201
	text := "Inventory created"

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

	var newItems models.InventoryItem
	err = json.Unmarshal(data, &newItems)
	if err != nil {
		statusCode = 400
		text = "Invalid JSON format"
		return
	}

	if err := h.Service.CreateInventory(newItems); err != nil {
		if errors.Is(err, cerrors.ErrExist) {
			statusCode = 409
			text = err.Error()
		}
		statusCode = 400
		text = err.Error()
		return
	}
}

func (h *Handler) GetInventoryItemId(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "Inventory got"

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

	data, err := h.Service.InventoryGetId(id)
	if err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = cerrors.ErrNotExist.Error()
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

func (h *Handler) UpDate(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "inventory updated"

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

	var newItems models.InventoryItem
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

	if err := h.Service.InventoryUpDate(id, newItems); err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = cerrors.ErrNotExist.Error()
		}

		statusCode = 400
		text = err.Error()
		return
	}
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	statusCode := 204
	text := "inventory deleted"

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

	if err := h.Service.DeleteInvent(id); err != nil {
		if errors.Is(err, cerrors.ErrNotExist) {
			statusCode = 404
			text = cerrors.ErrNotExist.Error()
			return
		}
		statusCode = 400
		text = err.Error()
		return
	}
}

func (h *Handler) GetLeftOvers(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из query
	sortBy := r.URL.Query().Get("sortBy")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Устанавливаем значения по умолчанию, если параметры не переданы
	if pageStr == "" {
		pageStr = "1" // Страница по умолчанию - 1
	}
	if pageSizeStr == "" {
		pageSizeStr = "10" // Количество товаров на странице по умолчанию - 10
	}

	// Преобразуем page и pageSize в int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		http.Error(w, "Invalid pageSize parameter", http.StatusBadRequest)
		return
	}

	// Логика для сортировки
	if sortBy != "price" && sortBy != "quantity" && sortBy != "" {
		http.Error(w, "Invalid sortBy parameter", http.StatusBadRequest)
		return
	}

	// Получаем остатки инвентаря с учетом параметров
	results, err := h.Service.GetLeftOvers(sortBy, page, pageSize)
	if err != nil {
		http.Error(w, "Failed to get leftovers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем результат в формате JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
