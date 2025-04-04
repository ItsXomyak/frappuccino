package internal

import (
	"database/sql"

	"frappuccino/internal/repo/invent"
	"frappuccino/internal/repo/menu"
	"frappuccino/internal/repo/order"
	"frappuccino/internal/repo/search"
)

type Container struct {
	MenuRepo      menu.MenuRepository
	InventoryRepo invent.Inventory
	OrderRepo     order.OrderRepository
	SearchRepo    search.SearchRepository
}

func New(path *sql.DB) *Container {
	return &Container{
		MenuRepo:      menu.New(path),
		InventoryRepo: invent.New(path),
		OrderRepo:     order.New(path),
		SearchRepo:    search.New(path),
	}
}
