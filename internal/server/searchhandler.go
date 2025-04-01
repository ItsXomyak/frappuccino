package server

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) FullTextSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")

	if query == "" {
		http.Error(w, "Missing required query parameter 'q'", http.StatusBadRequest)
		return
	}

	results, err := h.Service.SearchFullText(query, filter, minPrice, maxPrice)
	if err != nil {
		http.Error(w, "Failed to perform search: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetOrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	if period == "" {
		http.Error(w, "Missing required query parameter 'period'", http.StatusBadRequest)
		return
	}

	report, err := h.Service.GetOrderedItemsByPeriod(period, month, year)
	if err != nil {
		http.Error(w, "Failed to get ordered items by period: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(report)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
