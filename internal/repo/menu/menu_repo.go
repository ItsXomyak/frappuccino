package menu

import (
	"database/sql"
	"fmt"
	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
	"log/slog"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type MenuRepository interface {
	CreateMenuItem(item models.MenuItem) (*models.MenuItem, error)
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItemByID(id int) (*models.MenuItem, error)
	UpdateMenuItem(id int, item models.MenuItem) (*models.MenuItem, error)
	DeleteMenuItem(id int) error
	GetIngredientsByMenuItemID(menuItemID int) ([]models.MenuItemIngredient, error) // Новый метод
}

type menuRepository struct {
	db *sql.DB
}

func New(db *sql.DB) MenuRepository {
	return &menuRepository{
		db: db,
	}
}

func (r *menuRepository) CreateMenuItem(item models.MenuItem) (*models.MenuItem, error) {
	var id int
	err := r.db.QueryRow(`
        INSERT INTO menu_items (name, description, categories, allergens, price, available, size)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		item.Name, item.Description, pq.Array(item.Categories), pq.Array(item.Allergens), item.Price, item.Available, item.Size).
		Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu item: %v", err)
	}

	item.ID = id
	slog.Info("Menu item created", "item_id", id, "name", item.Name)
	return &item, nil
}

func (r *menuRepository) GetIngredientsByMenuItemID(menuItemID int) ([]models.MenuItemIngredient, error) {
	rows, err := r.db.Query(`
        SELECT ingredient_id, quantity
        FROM menu_item_ingredients
        WHERE menu_item_id = $1`, menuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %v", err)
	}
	defer rows.Close()

	var ingredients []models.MenuItemIngredient
	for rows.Next() {
		var ing models.MenuItemIngredient
		if err := rows.Scan(&ing.IngredientID, &ing.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %v", err)
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients, nil
}

func (r *menuRepository) GetAllMenuItems() ([]models.MenuItem, error) {
	rows, err := r.db.Query(`
        SELECT id, name, description, categories, allergens, price, available, size
        FROM menu_items`)
	if err != nil {
		return nil, fmt.Errorf("failed to query menu items: %v", err)
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, pq.Array(&item.Categories), pq.Array(&item.Allergens), &item.Price, &item.Available, &item.Size)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %v", err)
		}
		items = append(items, item)
	}

	slog.Info("GetAllMenuItems called", "count", len(items))
	return items, nil
}

func (r *menuRepository) GetMenuItemByID(id int) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.QueryRow(`
        SELECT id, name, description, categories, allergens, price, available, size
        FROM menu_items
        WHERE id = $1`, id).
		Scan(&item.ID, &item.Name, &item.Description, pq.Array(&item.Categories), pq.Array(&item.Allergens), &item.Price, &item.Available, &item.Size)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cerrors.ErrMenuItemNotFound
		}
		return nil, fmt.Errorf("failed to query menu item: %v", err)
	}

	slog.Info("Menu item found", "item_id", item.ID, "name", item.Name)
	return &item, nil
}

func (r *menuRepository) UpdateMenuItem(id int, item models.MenuItem) (*models.MenuItem, error) {
	result, err := r.db.Exec(`
        UPDATE menu_items 
        SET name = $1, description = $2, categories = $3, allergens = $4, price = $5, available = $6, size = $7
        WHERE id = $8`,
		item.Name, item.Description, pq.Array(item.Categories), pq.Array(item.Allergens), item.Price, item.Available, item.Size, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update menu item: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return nil, cerrors.ErrMenuItemNotFound
	}

	// Записываем историю изменения цены, если цена изменилась
	var oldPrice float64
	err = r.db.QueryRow(`SELECT price FROM menu_items WHERE id = $1`, id).Scan(&oldPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to get old price: %v", err)
	}
	if oldPrice != item.Price {
		_, err = r.db.Exec(`
            INSERT INTO price_history (menu_item_id, old_price, new_price, changed_at)
            VALUES ($1, $2, $3, NOW())`,
			id, oldPrice, item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to log price history: %v", err)
		}
	}

	item.ID = id
	slog.Info("Menu item updated", "item_id", id, "name", item.Name)
	return &item, nil
}

func (r *menuRepository) DeleteMenuItem(id int) error {
	result, err := r.db.Exec(`DELETE FROM menu_items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu item: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return cerrors.ErrMenuItemNotFound
	}

	slog.Info("Menu item deleted", "item_id", id)
	return nil
}
