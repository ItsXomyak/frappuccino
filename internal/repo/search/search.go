package search

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"frappuccino/internal/models"
)

type SearchRepository interface {
	SearchAll(query, filter, minPrice, maxPrice string) ([]models.Search, error)
	SearchOrders(query, filter, minPrice, maxPrice string) ([]*models.Search, error)
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
		whereConditions = append(whereConditions, "(name ILIKE '"+queryString+"' OR description ILIKE '"+queryString+"')")
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

	sqlQuery := "SELECT id, name, description, price FROM menu_items"
	whereConditions := []string{}

	if filter == "menu" || filter == "all" {
		whereConditions = append(whereConditions, "(name ILIKE '"+queryString+"' OR description ILIKE '"+queryString+"')")
	}

	if minPrice != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("price >= %s", minPrice))
	}

	if maxPrice != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("price <= %s", maxPrice))
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

func (r *searchRepository) SearchOrders(query, filter, minPrice, maxPrice string) ([]*models.Search, error) {
	var results []*models.Search
	queryString := "%" + query + "%"

	// Основной запрос для поиска заказов
	sqlQuery := `
        SELECT o.id, c.name AS customer_name, oi.order_details, o.total_amount
        FROM orders o
        JOIN customers c ON o.customer_id = c.id
        LEFT JOIN (
            SELECT oi.order_id, STRING_AGG(mi.name, ', ') AS order_details
            FROM order_items oi
            JOIN menu_items mi ON oi.menu_item_id = mi.id
            GROUP BY oi.order_id
        ) oi ON o.id = oi.order_id
    `

	whereConditions := []string{}

	// Условия фильтрации
	if filter == "orders" || filter == "all" {
		whereConditions = append(whereConditions, "(c.name ILIKE $1 OR oi.order_details ILIKE $1)")
	}

	// Фильтрация по цене
	if minPrice != "" {
		whereConditions = append(whereConditions, "o.total_amount >= "+minPrice)
	}

	if maxPrice != "" {
		whereConditions = append(whereConditions, "o.total_amount <= "+maxPrice)
	}

	// Добавление условий WHERE
	if len(whereConditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Выполнение запроса
	var rows *sql.Rows
	var err error

	if filter == "orders" || filter == "all" {
		rows, err = r.db.Query(sqlQuery, queryString)
	} else {
		rows, err = r.db.Query(sqlQuery)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	// Чтение строк результата
	for rows.Next() {
		result := &models.Search{}
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

	period = strings.ToLower(period)
	month = strings.ToLower(month)

	if period == "day" {
		if month == "" {
			return nil, fmt.Errorf("month is required when period is 'day'")
		}

		// Преобразуем месяц в числовое значение
		monthInt, err := time.Parse("January", strings.Title(month)) // "october" -> "October"
		if err != nil {
			return nil, fmt.Errorf("invalid month name: %w", err)
		}
		monthNumber := int(monthInt.Month())

		// Если год не передан, используем текущий год
		if year == "" {
			year = fmt.Sprintf("%d", time.Now().Year())
		}

		// Преобразуем год в целое число
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return nil, fmt.Errorf("invalid year: %w", err)
		}

		// Формируем запрос для дня
		queryStr = `
			SELECT EXTRACT(DAY FROM created_at)::integer AS day, COUNT(*) AS count
			FROM orders
			WHERE EXTRACT(MONTH FROM created_at) = $1
			AND EXTRACT(YEAR FROM created_at) = $2
			GROUP BY EXTRACT(DAY FROM created_at)
			ORDER BY day
		`

		// Выполняем запрос с параметрами
		rows, err := r.db.Query(queryStr, monthNumber, yearInt)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		defer rows.Close()

		// Чтение строк результата
		for rows.Next() {
			var day int
			var count int
			err := rows.Scan(&day, &count)
			if err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
			// Используем день как строку без ведущих нулей
			dayStr := fmt.Sprintf("%d", day)
			report = append(report, models.OrderedItemReport{
				Period: dayStr,
				Count:  count,
			})
		}

		// Проверка на ошибки после выполнения запроса
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error after scanning rows: %w", err)
		}
	}

	// Проверяем период 'month'
	if period == "month" {
		if year == "" {
			return nil, fmt.Errorf("year is required when period is 'month'")
		}

		// Преобразуем год в целое число
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return nil, fmt.Errorf("invalid year: %w", err)
		}

		// Формируем запрос для месяца
		queryStr = `
			SELECT LOWER(TO_CHAR(created_at, 'FMMonth')) AS month, COUNT(*) AS count
			FROM orders
			WHERE EXTRACT(YEAR FROM created_at) = $1
			GROUP BY EXTRACT(MONTH FROM created_at), TO_CHAR(created_at, 'FMMonth')
			ORDER BY EXTRACT(MONTH FROM created_at)
		`

		// Выполняем запрос с параметром года
		rows, err := r.db.Query(queryStr, yearInt)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		defer rows.Close()

		// Чтение строк результата
		for rows.Next() {
			var monthName string
			var count int
			err := rows.Scan(&monthName, &count)
			if err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
			report = append(report, models.OrderedItemReport{
				Period: monthName, // Уже в нижнем регистре из запроса
				Count:  count,
			})
		}

		// Проверка на ошибки после выполнения запроса
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error after scanning rows: %w", err)
		}
	}

	// Проверка на случай, если period некорректен
	if period != "day" && period != "month" {
		return nil, fmt.Errorf("invalid period: %s. Expected 'day' or 'month'", period)
	}

	return report, nil
}
