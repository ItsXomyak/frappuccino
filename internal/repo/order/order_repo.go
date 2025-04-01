package order

import (
	"database/sql"
	"fmt"
	"frappuccino/internal/models"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type OrderRepository interface {
	CreateOrder(data models.Order) error
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id int) (models.Order, error)
	UpdateOrder(id int, data models.Order) error
	DeleteOrder(id int) error
	CloseOrder(id int) error
	GetPopularItems() ([]models.PopularItem, error)
	GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error)
	GetCustomerNameByID(customerID int) (string, error)
	BatchProcessOrders(orders []models.Order) (*models.BatchOrderResponse, error)
}

type orderRepository struct {
	db *sql.DB
}

func New(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateOrder(data models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Вставка заказа в таблицу orders
	var orderID int
	err = tx.QueryRow(`
        INSERT INTO orders (customer_id, status, total_amount, payment_method, special_instructions, created_at)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		data.CustomerID, "open", data.TotalAmount, data.PaymentMethod, data.SpecialInstructions, time.Now()).
		Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	// Вставка элементов заказа и обновление инвентаря
	for _, item := range data.Items {
		// Получаем цену из menu_items
		var price float64
		err = tx.QueryRow(`SELECT price FROM menu_items WHERE id = $1`, item.MenuItemID).Scan(&price)
		if err != nil {
			return fmt.Errorf("failed to get menu item price: %v", err)
		}

		// Вставка элемента заказа
		var itemID int
		err = tx.QueryRow(`
            INSERT INTO order_items (order_id, menu_item_id, quantity, price, customizations)
            VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			orderID, item.MenuItemID, item.Quantity, price, item.Customizations).Scan(&itemID)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %v", err)
		}

		// Проверка и обновление инвентаря
		rows, err := tx.Query(`
            SELECT ingredient_id, quantity FROM menu_item_ingredients WHERE menu_item_id = $1`, item.MenuItemID)
		if err != nil {
			return fmt.Errorf("failed to get ingredients: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var ingredientID int
			var requiredQty float64
			if err := rows.Scan(&ingredientID, &requiredQty); err != nil {
				return fmt.Errorf("failed to scan ingredient: %v", err)
			}

			var stock float64
			err = tx.QueryRow(`SELECT stock FROM inventory WHERE id = $1 FOR UPDATE`, ingredientID).Scan(&stock)
			if err != nil {
				return fmt.Errorf("failed to check inventory: %v", err)
			}

			totalRequired := requiredQty * float64(item.Quantity)
			if stock < totalRequired {
				return fmt.Errorf("insufficient stock for ingredient %d", ingredientID)
			}

			_, err = tx.Exec(`
                UPDATE inventory SET stock = stock - $1 WHERE id = $2`,
				totalRequired, ingredientID)
			if err != nil {
				return fmt.Errorf("failed to update inventory: %v", err)
			}

			_, err = tx.Exec(`
                INSERT INTO inventory_transactions (ingredient_id, change_amount, transaction_type, occurred_at)
                VALUES ($1, $2, 'use', NOW())`,
				ingredientID, -totalRequired)
			if err != nil {
				return fmt.Errorf("failed to log inventory transaction: %v", err)
			}
		}
	}

	// Запись начального статуса в историю
	_, err = tx.Exec(`
        INSERT INTO order_status_history (order_id, status, changed_at)
        VALUES ($1, 'open', NOW())`, orderID)
	if err != nil {
		return fmt.Errorf("failed to log order status: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	data.ID = orderID
	return nil
}

func (r *orderRepository) GetCustomerNameByID(customerID int) (string, error) {
	var name string
	query := "SELECT name FROM customers WHERE id = $1"
	err := r.db.QueryRow(query, customerID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("customer not found with id: %d", customerID)
		}
		return "", fmt.Errorf("failed to query customer: %w", err)
	}
	return name, nil
}

func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	rows, err := r.db.Query(`
        SELECT o.id, o.customer_id, o.status, o.total_amount, o.payment_method, o.special_instructions, o.created_at, o.updated_at,
               oi.id AS item_id, oi.menu_item_id, oi.quantity, oi.price, oi.customizations
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id`)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	ordersMap := make(map[int]*models.Order)
	for rows.Next() {
		var o models.Order
		var item models.OrderItem
		var itemID sql.NullInt64

		err := rows.Scan(&o.ID, &o.CustomerID, &o.Status, &o.TotalAmount, &o.PaymentMethod, &o.SpecialInstructions, &o.CreatedAt, &o.UpdatedAt,
			&itemID, &item.MenuItemID, &item.Quantity, &item.Price, &item.Customizations)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		if _, exists := ordersMap[o.ID]; !exists {
			o.Items = []models.OrderItem{}
			ordersMap[o.ID] = &o
		}
		if itemID.Valid {
			item.ID = int(itemID.Int64)
			ordersMap[o.ID].Items = append(ordersMap[o.ID].Items, item)
		}
	}

	var orders []models.Order
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (r *orderRepository) GetOrderByID(id int) (models.Order, error) {
	rows, err := r.db.Query(`
        SELECT o.id, o.customer_id, o.status, o.total_amount, o.payment_method, o.special_instructions, o.created_at, o.updated_at,
               oi.id AS item_id, oi.menu_item_id, oi.quantity, oi.price, oi.customizations
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        WHERE o.id = $1`, id)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to query order: %v", err)
	}
	defer rows.Close()

	var order models.Order
	first := true
	for rows.Next() {
		var item models.OrderItem
		var itemID sql.NullInt64

		if first {
			err = rows.Scan(&order.ID, &order.CustomerID, &order.Status, &order.TotalAmount, &order.PaymentMethod, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
				&itemID, &item.MenuItemID, &item.Quantity, &item.Price, &item.Customizations)
			first = false
		} else {
			err = rows.Scan(new(int), new(int), new(string), new(float64), new(string), new(string), new(time.Time), new(time.Time),
				&itemID, &item.MenuItemID, &item.Quantity, &item.Price, &item.Customizations)
		}
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to scan order: %v", err)
		}
		if itemID.Valid {
			item.ID = int(itemID.Int64)
			order.Items = append(order.Items, item)
		}
	}
	if first {
		return models.Order{}, fmt.Errorf("order with ID %d not found", id)
	}
	return order, nil
}

func (r *orderRepository) UpdateOrder(id int, data models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        UPDATE orders 
        SET customer_id = $1, status = $2, total_amount = $3, payment_method = $4, special_instructions = $5
        WHERE id = $6`,
		data.CustomerID, "open", data.TotalAmount, data.PaymentMethod, data.SpecialInstructions, id)
	if err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Удаляем старые элементы и добавляем новые
	_, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete old items: %v", err)
	}

	for _, item := range data.Items {
		_, err = tx.Exec(`
            INSERT INTO order_items (order_id, menu_item_id, quantity, price, customizations)
            VALUES ($1, $2, $3, $4, $5)`,
			id, item.MenuItemID, item.Quantity, item.Price, item.Customizations)
		if err != nil {
			return fmt.Errorf("failed to insert new item: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func (r *orderRepository) DeleteOrder(id int) error {
	_, err := r.db.Exec(`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}
	return nil
}

func (r *orderRepository) CloseOrder(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE orders SET status = 'closed', updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to close order: %v", err)
	}

	_, err = tx.Exec(`
        INSERT INTO order_status_history (order_id, status, changed_at)
        VALUES ($1, 'closed', NOW())`, id)
	if err != nil {
		return fmt.Errorf("failed to log status: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func (r *orderRepository) GetPopularItems() ([]models.PopularItem, error) {
	rows, err := r.db.Query(`
        SELECT oi.menu_item_id, mi.name, SUM(oi.quantity) AS popularity
        FROM order_items oi
        JOIN menu_items mi ON oi.menu_item_id = mi.id
        GROUP BY oi.menu_item_id, mi.name
        ORDER BY popularity DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query popular items: %v", err)
	}
	defer rows.Close()

	var popularItems []models.PopularItem
	for rows.Next() {
		var item models.PopularItem
		if err := rows.Scan(&item.MenuItemID, &item.Name, &item.Popularity); err != nil {
			return nil, fmt.Errorf("failed to scan popular item: %v", err)
		}
		popularItems = append(popularItems, item)
	}
	return popularItems, nil
}

func (r *orderRepository) GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error) {
	queryStr := "SELECT item_name, COUNT(*) FROM orders"

	// Создаем условия фильтрации по дате
	var conditions []string
	if startDate != "" {
		conditions = append(conditions, fmt.Sprintf("order_date >= '%s'", startDate))
	}
	if endDate != "" {
		conditions = append(conditions, fmt.Sprintf("order_date <= '%s'", endDate))
	}

	if len(conditions) > 0 {
		queryStr += " WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr += " GROUP BY item_name"

	rows, err := r.db.Query(queryStr)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	results := make(map[string]int)

	for rows.Next() {
		var itemName string
		var count int
		err := rows.Scan(&itemName, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results[itemName] = count
	}

	// Проверка на ошибки после завершения перебора
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return results, nil
}

func (r *orderRepository) BatchProcessOrders(orders []models.Order) (*models.BatchOrderResponse, error) {
	// Начинаем транзакцию для обеспечения целостности данных
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // В случае ошибки откатим транзакцию

	var processedOrders []models.ProcessedOrder
	var totalRevenue float64
	var inventoryUpdates []models.InventoryUpdate

	// Обрабатываем каждый заказ
	for _, order := range orders {
		// Получаем имя клиента
		customerName, err := r.GetCustomerNameByID(order.CustomerID)
		if err != nil {
			// Если ошибка при получении имени, отклоняем заказ
			processedOrders = append(processedOrders, models.ProcessedOrder{
				CustomerName: customerName,
				Status:       "rejected",
				Reason:       err.Error(),
			})
			continue
		}

		// Проверка инвентаря (псевдокод, нужно добавить реальную проверку)
		for _, item := range order.Items {
			// Пример: проверка, что у нас есть достаточно товара в инвентаре
			var inventoryQuantity int
			err := tx.QueryRow("SELECT quantity FROM inventory WHERE item_id = $1", item.MenuItemID).Scan(&inventoryQuantity)
			if err != nil {
				// Если товара недостаточно
				processedOrders = append(processedOrders, models.ProcessedOrder{
					CustomerName: customerName,
					Status:       "rejected",
					Reason:       fmt.Sprintf("insufficient inventory for item ID %d", item.MenuItemID),
				})
				continue
			}

			// Проверяем, есть ли достаточное количество товара
			if inventoryQuantity < item.Quantity {
				processedOrders = append(processedOrders, models.ProcessedOrder{
					CustomerName: customerName,
					Status:       "rejected",
					Reason:       fmt.Sprintf("insufficient inventory for item %d", item.MenuItemID),
				})
				continue
			}

			// Обновляем количество в инвентаре
			_, err = tx.Exec("UPDATE inventory SET quantity = quantity - $1 WHERE item_id = $2", item.Quantity, item.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to update inventory: %w", err)
			}

			// Добавляем изменения в инвентарь в список
			inventoryUpdates = append(inventoryUpdates, models.InventoryUpdate{
				IngredientID: item.MenuItemID,
				Name:         fmt.Sprintf("Item %d", item.MenuItemID), // Это стоит заменить на реальное название товара
				QuantityUsed: item.Quantity,
				Remaining:    inventoryQuantity - item.Quantity,
			})
		}

		// Успешная обработка заказа
		processedOrders = append(processedOrders, models.ProcessedOrder{
			CustomerName: customerName,
			Status:       "accepted",
			Total:        order.TotalAmount,
		})

		// Подсчитываем общий доход
		totalRevenue += order.TotalAmount
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	summary := models.BatchOrderSummary{
		TotalOrders:      len(orders),
		Accepted:         countAcceptedOrders(processedOrders),
		Rejected:         len(orders) - countAcceptedOrders(processedOrders),
		TotalRevenue:     totalRevenue,
		InventoryUpdates: inventoryUpdates,
	}

	return &models.BatchOrderResponse{
		ProcessedOrders: processedOrders,
		Summary:         summary,
	}, nil
}

func countAcceptedOrders(orders []models.ProcessedOrder) int {
	acceptedCount := 0
	for _, order := range orders {
		if order.Status == "accepted" {
			acceptedCount++
		}
	}
	return acceptedCount
}
