package search

import (
	"database/sql"
	"fmt"
	"frappuccino/internal/models"
	"strconv"
	"strings"
	"time"
)

type SearchRepository interface {
	SearchAll(query, filter, minPrice, maxPrice string) ([]models.Search, error)
	SearchOrders(query, filter, minPrice, maxPrice string) ([]models.Search, error)
	SearchMenuItems(query, filter, minPrice, maxPrice string) ([]models.Search, error)
	GetOrderedItemsByPeriod(period, month, year string) ([]models.OrderedItemReport, error)
}

type searchRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *searchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchAll(query, filter, minPrice, maxPrice string) ([]models.Search, error) {
	var results []models.Search
	queryString := "%" + query + "%"

	sqlQuery := "SELECT id, name, description, price FROM items"
	whereConditions := []string{}

	// условия поиска
	if filter == "orders" || filter == "all" {
		whereConditions = append(whereConditions, "(customer_name ILIKE '"+queryString+"' OR order_details ILIKE '"+queryString+"')")
	}

	if filter == "menu" || filter == "all" {
		whereConditions = append(whereConditions, "(item_name ILIKE '"+queryString+"' OR item_description ILIKE '"+queryString+"')")
	}

	// фильтрации по цене
	if minPrice != "" {
		minPriceFloat, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid minPrice value: %w", err)
		}
		whereConditions = append(whereConditions, fmt.Sprintf("price >= %f", minPriceFloat))
	}

	if maxPrice != "" {
		maxPriceFloat, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid maxPrice value: %w", err)
		}
		whereConditions = append(whereConditions, fmt.Sprintf("price <= %f", maxPriceFloat))
	}

	if len(whereConditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereConditions, " OR ")
	}

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result models.Search
		err := rows.Scan(&result.ID, &result.Name, &result.Description, &result.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *searchRepository) SearchMenuItems(query, filter, minPrice, maxPrice string) ([]models.Search, error) {
	var results []models.Search
	queryString := "%" + query + "%"

	// Подготавливаем базовые части запроса
	sqlQuery := "SELECT id, item_name, item_description, price FROM menu_items"
	whereConditions := []string{}

	// Формируем условия поиска
	if filter == "menu" || filter == "all" {
		whereConditions = append(whereConditions, "(item_name ILIKE '"+queryString+"' OR item_description ILIKE '"+queryString+"')")
	}

	// Фильтрация по минимальной цене
	if minPrice != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("price >= %s", minPrice))
	}

	// Фильтрация по максимальной цене
	if maxPrice != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("price <= %s", maxPrice))
	}

	// Добавляем WHERE условия к запросу, если они есть
	if len(whereConditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereConditions, " OR ")
	}

	// Выполняем запрос без использования args...
	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	// Чтение строк результата
	for rows.Next() {
		var result models.Search
		err := rows.Scan(&result.ID, &result.Name, &result.Description, &result.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *searchRepository) SearchOrders(query, filter, minPrice, maxPrice string) ([]models.Search, error) {
	var results []models.Search
	queryString := "%" + query + "%"

	// Подготавливаем базовые части запроса
	sqlQuery := "SELECT id, customer_name, order_details, total FROM orders"
	whereConditions := []string{}

	if filter == "orders" || filter == "all" {
		whereConditions = append(whereConditions, "(customer_name ILIKE '"+queryString+"' OR order_details ILIKE '"+queryString+"')")
	}

	if minPrice != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("total >= %s", minPrice))
	}

	if maxPrice != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("total <= %s", maxPrice))
	}

	// Добавляем WHERE условия к запросу, если они есть
	if len(whereConditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereConditions, " OR ")
	}

	// Выполняем запрос без использования args...
	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	// Чтение строк результата
	for rows.Next() {
		var result models.Search
		err := rows.Scan(&result.ID, &result.Name, &result.Description, &result.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	// Проверка на ошибки после выполнения запроса
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *searchRepository) GetOrderedItemsByPeriod(period, month, year string) ([]models.OrderedItemReport, error) {
	var queryStr string
	var report []models.OrderedItemReport

	if period == "day" {
		if month == "" {
			return nil, fmt.Errorf("month is required when period is 'day'")
		}

		// Преобразуем название месяца в число
		monthInt, err := time.Parse("January", month)
		if err != nil {
			return nil, fmt.Errorf("invalid month name: %w", err)
		}
		monthNumber := int(monthInt.Month())

		queryStr = fmt.Sprintf(`
			SELECT EXTRACT(DAY FROM order_date) AS day, COUNT(*)
			FROM orders
			WHERE EXTRACT(MONTH FROM order_date) = %d
			AND EXTRACT(YEAR FROM order_date) = %s
			GROUP BY day ORDER BY day
		`, monthNumber, year)
	}

	if period == "month" {
		if year == "" {
			return nil, fmt.Errorf("year is required when period is 'month'")
		}

		queryStr = fmt.Sprintf(`
			SELECT TO_CHAR(order_date, 'Month') AS month, COUNT(*)
			FROM orders
			WHERE EXTRACT(YEAR FROM order_date) = %s
			GROUP BY month ORDER BY EXTRACT(MONTH FROM order_date)
		`, year)
	}

	// Выполняем запрос к базе данных без args...
	rows, err := r.db.Query(queryStr)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Чтение строк результата
	for rows.Next() {
		var monthOrDay string
		var count int
		err := rows.Scan(&monthOrDay, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		report = append(report, models.OrderedItemReport{
			Period: monthOrDay,
			Count:  count,
		})
	}

	// Проверка на ошибки после выполнения запроса
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return report, nil
}
