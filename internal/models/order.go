package models

import "time"

type Order struct {
	ID                  int         `json:"id"`
	CustomerID          int         `json:"customer_id"`
	Items               []OrderItem `json:"items"`
	Status              string      `json:"status"`
	TotalAmount         float64     `json:"total_amount"`
	PaymentMethod       string      `json:"payment_method"`
	SpecialInstructions string      `json:"special_instructions"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID             int     `json:"id"`
	OrderID        int     `json:"order_id"`
	MenuItemID     int     `json:"menu_item_id"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	Customizations string  `json:"customizations"`
}

type ProcessedOrder struct {
	OrderID      int     `json:"order_id,omitempty"`
	CustomerName string  `json:"customer_name"`
	Status       string  `json:"status"`
	Reason       string  `json:"reason,omitempty"`
	Total        float64 `json:"total,omitempty"`
}

type InventoryUpdate struct {
	IngredientID int    `json:"ingredient_id"`
	Name         string `json:"name"`
	QuantityUsed int    `json:"quantity_used"`
	Remaining    int    `json:"remaining"`
}

type BatchOrderSummary struct {
	TotalOrders      int               `json:"total_orders"`
	Accepted         int               `json:"accepted"`
	Rejected         int               `json:"rejected"`
	TotalRevenue     float64           `json:"total_revenue"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}

type BatchOrderResponse struct {
	ProcessedOrders []ProcessedOrder  `json:"processed_orders"`
	Summary         BatchOrderSummary `json:"summary"`
}
