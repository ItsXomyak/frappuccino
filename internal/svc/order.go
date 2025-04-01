package svc

import (
	"fmt"
	"frappuccino/helper"
	"frappuccino/internal/models"
	"time"
)

func (s *svc) OrderCreate(data models.Order) error {
	if data.CustomerID <= 0 {
		s.Log.Error("Invalid customer ID", "customer_id", data.CustomerID)
		return fmt.Errorf("invalid customer ID")
	}

	dataMenu, err := s.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu", "error", err.Error())
		return err
	}

	if err := helper.CheckForOrders(data, dataMenu); err != nil {
		s.Log.Error("Order validation failed", "customer_id", data.CustomerID, "error", err.Error())
		return err
	}

	data.Status = "open"
	data.CreatedAt = time.Now()
	if err := s.Repo.OrderRepo.CreateOrder(data); err != nil {
		s.Log.Error("Failed to create order", "customer_id", data.CustomerID, "error", err.Error())
		return err
	}

	s.Log.Info("Successfully created order", "id", data.ID)
	return nil
}

func (s *svc) Get() ([]models.Order, error) {
	data, err := s.Repo.OrderRepo.GetAllOrders()
	if err != nil {
		s.Log.Error("Failed to retrieve orders", "error", err.Error())
		return nil, err
	}

	s.Log.Info("Successfully retrieved all orders", "count", len(data))
	return data, nil
}

func (s *svc) GetId(id int) (models.Order, error) {
	dataId, err := s.Repo.OrderRepo.GetOrderByID(id)
	if err != nil {
		s.Log.Error("Failed to retrieve order by ID", "id", id, "error", err.Error())
		return models.Order{}, err
	}

	s.Log.Info("Successfully retrieved order", "id", id)
	return dataId, nil
}

func (s *svc) Update(id int, data models.Order) error {
	if data.CustomerID <= 0 {
		s.Log.Error("Invalid customer ID", "customer_id", data.CustomerID)
		return fmt.Errorf("invalid customer ID")
	}

	dataMenu, err := s.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu", "error", err.Error())
		return err
	}

	if err := helper.CheckForOrders(data, dataMenu); err != nil {
		s.Log.Error("Order validation failed", "error", err.Error())
		return err
	}

	if err := s.Repo.OrderRepo.UpdateOrder(id, data); err != nil {
		s.Log.Error("Failed to update order", "id", id, "error", err.Error())
		return err
	}

	s.Log.Info("Successfully updated order", "id", id)
	return nil
}

func (s *svc) RemoveOrder(id int) error {
	if err := s.Repo.OrderRepo.DeleteOrder(id); err != nil {
		s.Log.Error("Failed to delete order", "id", id, "error", err.Error())
		return err
	}

	s.Log.Info("Successfully deleted order", "id", id)
	return nil
}

func (s *svc) CloseOrder(id int) error {
	order, err := s.Repo.OrderRepo.GetOrderByID(id)
	if err != nil {
		s.Log.Error("Failed to retrieve order by ID", "id", id, "error", err.Error())
		return err
	}

	if order.Status == "closed" {
		s.Log.Error("Operation not allowed: status is closed", "id", id)
		return fmt.Errorf("operation not allowed: status is closed")
	}

	menuAll, err := s.Repo.MenuRepo.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu items", "error", err.Error())
		return err
	}

	menuMap := make(map[int]models.MenuItem)
	for _, item := range menuAll {
		menuMap[item.ID] = item
	}

	inventoryAll, err := s.Repo.InventoryRepo.GetInventory()
	if err != nil {
		s.Log.Error("Failed to retrieve inventory", "error", err.Error())
		return err
	}

	ingredientMap := make(map[int]models.InventoryItem)
	for _, item := range inventoryAll {
		ingredientMap[item.ID] = item
	}

	for _, orderItem := range order.Items {
		if menuItem, exists := menuMap[orderItem.MenuItemID]; exists {
			ingredients, err := s.Repo.MenuRepo.GetIngredientsByMenuItemID(menuItem.ID)
			if err != nil {
				s.Log.Error("Failed to get ingredients", "menu_item_id", menuItem.ID, "error", err.Error())
				return err
			}

			for _, ingredient := range ingredients {
				requiredQty := ingredient.Quantity * float64(orderItem.Quantity)
				if invItem, exists := ingredientMap[ingredient.IngredientID]; exists {
					if requiredQty > invItem.Stock {
						s.Log.Error("Not enough quantity for ingredient", "ingredient_id", ingredient.IngredientID, "required", requiredQty, "available", invItem.Stock)
						return fmt.Errorf("not enough quantity for ingredient ID %d: required %v, available %v", ingredient.IngredientID, requiredQty, invItem.Stock)
					}
					invItem.Stock -= requiredQty
					ingredientMap[ingredient.IngredientID] = invItem
				}
			}
		}
	}

	for ingredientID, item := range ingredientMap {
		if err := s.Repo.InventoryRepo.PutInventory(ingredientID, item); err != nil {
			s.Log.Error("Failed to update inventory", "ingredient_id", ingredientID, "error", err.Error())
			return err
		}
	}

	if err := s.Repo.OrderRepo.CloseOrder(id); err != nil {
		s.Log.Error("Failed to close order", "id", id, "error", err.Error())
		return err
	}

	s.Log.Info("Successfully closed order", "id", id)
	return nil
}

func (s *svc) GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error) {
	result, err := s.Repo.OrderRepo.GetNumberOfOrderedItems(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered items: %w", err)
	}

	return result, nil
}

func (s *svc) BatchProcessOrders(orders []models.Order) (*models.BatchOrderResponse, error) {
	return s.Repo.OrderRepo.BatchProcessOrders(orders)

}
