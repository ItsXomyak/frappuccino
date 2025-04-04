package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) GetTotalSaleS(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	var text string

	totalSales, err := h.Service.GetTotalSales()
	if err != nil {
		statusCode = 400
		text = "Failed to get total sales"
		Respond(w, statusCode, text) // Ensure that the response is sent immediately in case of error
		return
	}

	// Check if total sales is 0 and handle it accordingly
	if totalSales == 0 {
		text = fmt.Sprintf("No sales recorded for the closed orders")
		Respond(w, statusCode, text)
		return
	}

	// Return the actual total sales if valid
	response := map[string]float64{
		"total_sales": totalSales,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		statusCode = 400
		text = "Failed to encode response"
		Respond(w, statusCode, text)
	}
}

func (h *Handler) GetPopMenuItems(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "popular menu"

	defer func() {
		Respond(w, statusCode, text)
	}()

	popularItems, err := h.Service.GetPopularItems()
	if err != nil {
		statusCode = 400
		text = "Failed to get popular menu items"
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(popularItems); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetExpensiveMenuItem(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	text := "popular menu"

	defer func() {
		Respond(w, statusCode, text)
	}()

	data, err := h.Service.GetExpensiveMenuItem()
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
