package svc

import (
	"frappuccino/internal/models"
	repo "frappuccino/internal/repo"
	"log/slog"
	"os"
)

type Svc interface {
	CreateInventory(data models.InventoryItem) error
	InventoriesGet() ([]models.InventoryItem, error)
	InventoryGetId(id int) (models.InventoryItem, error)
	InventoryUpDate(id int, data models.InventoryItem) error
	DeleteInvent(id int) error
	CreateMenuItem(item models.MenuItem) (*models.MenuItem, error)
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItemByID(id int) (*models.MenuItem, error)
	UpdateMenuItem(id int, item models.MenuItem) (*models.MenuItem, error)
	DeleteMenuItem(id int) error
	OrderCreate(data models.Order) error
	Get() ([]models.Order, error)
	GetId(id int) (models.Order, error)
	RemoveOrder(id int) error
	CloseOrder(id int) error
	Update(id int, data models.Order) error
	GetPopularItems() ([]models.PopularItem, error)
	GetTotalSales() (float64, error)
	GetExpensiveMenuItem() (models.MenuItem, error)
	SearchFullText(query, filter, minPrice, maxPrice string) (*models.SearchResponse, error)
	GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error)
	GetOrderedItemsByPeriod(period, month, year string) ([]models.OrderedItemReport, error)
	BatchProcessOrders(orders []models.Order) (*models.BatchOrderResponse, error)
	GetLeftOvers(sortBy string, page int, pageSize int) (*models.InventoryResponse, error)
}

type svc struct {
	Repo *repo.Container
	Log  *slog.Logger
}

func NewSvc(r *repo.Container) Svc {
	loger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &svc{
		Repo: r,
		Log:  loger,
	}
}
