package models

type InventoryItem struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Stock            float64 `json:"stock"`
	Unit             string  `json:"unit"`
	ReorderThreshold float64 `json:"reorder_threshold"`
	Price            float64 `json:"price"`
}

type InventoryResponse struct {
	CurrentPage int             `json:"currentPage"`
	HasNextPage bool            `json:"hasNextPage"`
	PageSize    int             `json:"pageSize"`
	TotalPages  int             `json:"totalPages"`
	Data        []InventoryItem `json:"data"`
}
