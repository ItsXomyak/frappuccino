package helper

import (
	"fmt"

	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
)

func CheckForOrders(order models.Order, menuItems []models.MenuItem) error {
	if order.ID != 0 { // Проверяем, что ID не задан для нового заказа
		return fmt.Errorf("order validation failed: order ID is already set: %d", order.ID)
	}

	if order.CustomerID <= 0 { // Проверяем, что CustomerID валиден
		return fmt.Errorf("please provide a valid customer ID for the order")
	}

	if len(order.Items) == 0 { // Проверяем, что заказ содержит элементы
		return fmt.Errorf("order must contain at least one item")
	}

	menuMap := make(map[int]models.MenuItem)
	for _, item := range menuItems {
		menuMap[item.ID] = item
	}

	for _, orderItem := range order.Items {
		if orderItem.MenuItemID == 0 { // Проверяем MenuItemID
			return fmt.Errorf("please provide a menu item ID for one of the items in the order")
		}
		if _, exists := menuMap[orderItem.MenuItemID]; !exists {
			return fmt.Errorf("menu item with ID %d not found", orderItem.MenuItemID)
		}
		if orderItem.Quantity <= 0 {
			return fmt.Errorf("please specify a quantity greater than zero for the item with menu item ID %d", orderItem.MenuItemID)
		}
		if orderItem.Price < 0 { // Проверяем, что цена не отрицательная
			return fmt.Errorf("price for item with menu item ID %d cannot be negative", orderItem.MenuItemID)
		}
	}

	if order.Status != "" && order.Status != "open" { // Проверяем, что статус либо пустой, либо "open" для нового заказа
		return fmt.Errorf("new order must have status 'open' or be unset, got: %s", order.Status)
	}

	if order.TotalAmount < 0 { // Проверяем, что сумма не отрицательная
		return fmt.Errorf("total amount cannot be negative, got: %f", order.TotalAmount)
	}

	return nil
}

func CheckCustomerExists(orders []models.Order, customerID int) error { // Изменяем на customerID
	if orders == nil {
		return fmt.Errorf("orders list is nil")
	}

	for _, order := range orders {
		if order.CustomerID == customerID {
			return nil
		}
	}
	return cerrors.ErrNotExist
}

func CheckCustomerExistsId(orders []models.Order, id int) error {
	if orders == nil {
		return fmt.Errorf("orders list is nil")
	}
	for _, order := range orders {
		if order.ID == id {
			return nil
		}
	}
	return cerrors.ErrNotExist
}

func CheckStatus(orders []models.Order, id int) error {
	if orders == nil {
		return fmt.Errorf("orders list is nil")
	}
	for _, order := range orders {
		if order.ID == id {
			if order.Status == "closed" {
				return fmt.Errorf("order with ID %d is closed", id)
			}
			return nil
		}
	}
	return fmt.Errorf("order with ID %d not found", id)
}
