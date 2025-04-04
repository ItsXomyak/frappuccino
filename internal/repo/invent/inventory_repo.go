package invent

import (
	"database/sql"
	"fmt"

	"frappuccino/internal/models"

	_ "github.com/lib/pq"
)

type Inventory interface {
	CreateInventory(data models.InventoryItem) error
	GetInventoryId(id int) (models.InventoryItem, error)
	GetInventory() ([]models.InventoryItem, error)
	PutInventory(id int, data models.InventoryItem) error
	DeleteInvent(id int) error
	GetByNameAndUnit(name, unit string) (models.InventoryItem, error)
	GetLeftOversWithPagination(sortBy string, page int, pageSize int) ([]models.InventoryItem, error)
	CountTotalInventoryItems() (int, error)
}

type inventory struct {
	db *sql.DB
}

func New(db *sql.DB) Inventory {
	return &inventory{
		db: db,
	}
}

func (i *inventory) GetByNameAndUnit(name, unit string) (models.InventoryItem, error) {
	var item models.InventoryItem
	err := i.db.QueryRow(`
        SELECT id, name, stock, unit, reorder_threshold
        FROM inventory
        WHERE name = $1 AND unit = $2`, name, unit).
		Scan(&item.ID, &item.Name, &item.Stock, &item.Unit, &item.ReorderThreshold)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.InventoryItem{}, fmt.Errorf("item with name %s and unit %s not found", name, unit)
		}
		return models.InventoryItem{}, fmt.Errorf("failed to query inventory item: %v", err)
	}
	return item, nil
}

func (i *inventory) CreateInventory(data models.InventoryItem) error {
	var id int
	err := i.db.QueryRow(`
        INSERT INTO inventory (name, stock, unit, reorder_threshold, price)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		data.Name, data.Stock, data.Unit, data.ReorderThreshold, data.Price).
		Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create inventory item: %v", err)
	}
	data.ID = id
	return nil
}

func (i *inventory) GetInventoryId(id int) (models.InventoryItem, error) {
	var item models.InventoryItem
	err := i.db.QueryRow(`
        SELECT id, name, stock, unit, reorder_threshold
        FROM inventory
        WHERE id = $1`, id).
		Scan(&item.ID, &item.Name, &item.Stock, &item.Unit, &item.ReorderThreshold)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.InventoryItem{}, fmt.Errorf("item with ID %d not found", id)
		}
		return models.InventoryItem{}, fmt.Errorf("failed to query inventory item: %v", err)
	}
	return item, nil
}

func (i *inventory) GetInventory() ([]models.InventoryItem, error) {
	rows, err := i.db.Query(`
        SELECT id, name, stock, unit, reorder_threshold
        FROM inventory`)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory: %v", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(&item.ID, &item.Name, &item.Stock, &item.Unit, &item.ReorderThreshold)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (i *inventory) PutInventory(id int, upDate models.InventoryItem) error {
	result, err := i.db.Exec(`
        UPDATE inventory 
        SET name = $1, stock = $2, unit = $3, reorder_threshold = $4
        WHERE id = $5`,
		upDate.Name, upDate.Stock, upDate.Unit, upDate.ReorderThreshold, id)
	if err != nil {
		return fmt.Errorf("failed to update inventory item: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %d not found", id)
	}
	return nil
}

func (i *inventory) DeleteInvent(id int) error {
	// Проверка зависимостей в таблице menu_item_ingredients
	var count int
	err := i.db.QueryRow(`SELECT COUNT(*) FROM menu_item_ingredients WHERE ingredient_id = $1`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for related menu item ingredients: %v", err)
	}

	// Если зависимости найдены, удаляем их
	if count > 0 {
		_, err := i.db.Exec(`DELETE FROM menu_item_ingredients WHERE ingredient_id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete related menu item ingredients: %v", err)
		}
	}

	// Проверка зависимостей в таблице inventory_transactions
	err = i.db.QueryRow(`SELECT COUNT(*) FROM inventory_transactions WHERE ingredient_id = $1`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for related inventory transactions: %v", err)
	}

	// Если зависимости найдены, удаляем их
	if count > 0 {
		_, err := i.db.Exec(`DELETE FROM inventory_transactions WHERE ingredient_id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete related inventory transactions: %v", err)
		}
	}

	// Удаление элемента из таблицы inventory
	result, err := i.db.Exec(`DELETE FROM inventory WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory item: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %d not found", id)
	}

	return nil
}

func (r *inventory) GetLeftOversWithPagination(sortBy string, page int, pageSize int) ([]models.InventoryItem, error) {
	// Проверяем параметры на валидность
	if page < 1 {
		return nil, fmt.Errorf("invalid page parameter")
	}
	if pageSize < 1 {
		return nil, fmt.Errorf("invalid pageSize parameter")
	}

	var query string

	// Формируем запрос в зависимости от параметра сортировки
	switch sortBy {
	case "price":
		query = "SELECT id, name, stock, unit, reorder_threshold, price FROM inventory ORDER BY price LIMIT $1 OFFSET $2"
	case "quantity":
		query = "SELECT id, name, stock, unit, reorder_threshold, price FROM inventory ORDER BY stock LIMIT $1 OFFSET $2"
	default:
		query = "SELECT id, name, stock, unit, reorder_threshold, price FROM inventory ORDER BY name LIMIT $1 OFFSET $2"
	}

	// Рассчитываем OFFSET для пагинации
	offset := (page - 1) * pageSize

	// Выполняем запрос с пагинацией
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	// Сканы для каждой строки результата
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Stock, &item.Unit, &item.ReorderThreshold, &item.Price); err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, item)
	}

	// Проверка на ошибки во время итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %w", err)
	}

	return items, nil
}

func (r *inventory) CountTotalInventoryItems() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM inventory").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count inventory items: %w", err)
	}
	return count, nil
}
