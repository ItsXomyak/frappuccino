package models

type Search struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Relevance   float64 `json:"relevance"`
}

type SearchResponse struct {
	MenuItems    []Search `json:"menu_items"`
	Orders       []Search `json:"orders"`
	TotalMatches int      `json:"total_matches"`
}

type OrderedItemReport struct {
	Period string `json:"period"` 
	Count  int    `json:"count"`  
}
