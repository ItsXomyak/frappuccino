package models

type MenuItem struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Categories  []string             `json:"categories"`
	Allergens   []string             `json:"allergens"`
	Price       float64              `json:"price"`
	Available   bool                 `json:"available"`
	Size        string               `json:"size"`
	Ingredients []MenuItemIngredient `json:"ingredients,omitempty"`
}

type PopularItem struct {
	MenuItemID int    `json:"menu_item_id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

type MenuItemIngredient struct {
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}
