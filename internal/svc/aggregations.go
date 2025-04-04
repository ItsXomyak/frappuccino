package svc

import (
	"fmt"

	"frappuccino/internal/models"
)

func (s *svc) GetTotalSales() (float64, error) {
	orders, err := s.Repo.OrderRepo.GetAllOrders()
	if err != nil {
		s.Log.Error("Failed to retrieve orders", "error", err.Error())
		return 0, err
	}

	menuAll, err := s.Repo.MenuRepo.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu items", "error", err.Error())
		return 0, err
	}

	menuMap := make(map[int]models.MenuItem)
	for _, item := range menuAll {
		menuMap[item.ID] = item
	}

	total := 0.0
	for _, order := range orders {
		if order.Status != "closed" {
			continue
		}
		for _, item := range order.Items {
			if menuItem, exists := menuMap[item.MenuItemID]; exists {
				total += menuItem.Price * float64(item.Quantity)
			}
		}
	}

	if total == 0 {
		s.Log.Info("No closed orders found, total sales: 0", "total", total)
		return 0, nil // You can choose to return a meaningful message or a custom error
	}

	s.Log.Info("Successfully calculated total sales", "total", total)
	return total, nil
}

func (s *svc) GetPopularItems() ([]models.PopularItem, error) {
	popularItems, err := s.Repo.OrderRepo.GetPopularItems()
	if err != nil {
		s.Log.Error("Failed to retrieve popular items", "error", err.Error())
		return nil, err
	}

	s.Log.Info("Successfully retrieved popular items", "count", len(popularItems))
	return popularItems, nil
}

func (s *svc) GetExpensivMenuItems() (models.MenuItem, error) {
	menuAll, err := s.Repo.MenuRepo.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu items", "error", err.Error())
		return models.MenuItem{}, err
	}

	if len(menuAll) == 0 {
		s.Log.Error("Menu is empty")
		return models.MenuItem{}, fmt.Errorf("menu is empty")
	}

	var expensiveItem models.MenuItem

	for _, item := range menuAll {
		if item.Price > expensiveItem.Price {
			expensiveItem = item
		}
	}

	return expensiveItem, nil
}

func (s *svc) GetExpensiveMenuItem() (models.MenuItem, error) {
	menuAll, err := s.Repo.MenuRepo.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu items", "error", err.Error())
		return models.MenuItem{}, err
	}
	if len(menuAll) == 0 {
		s.Log.Info("Menu is empty")
		return models.MenuItem{}, fmt.Errorf("menu is empty")
	}
	expensiveItem := menuAll[0]
	for _, item := range menuAll[1:] {
		if item.Price > expensiveItem.Price {
			expensiveItem = item
		}
	}
	s.Log.Info("Successfully found expensive menu item", "id", expensiveItem.ID, "price", expensiveItem.Price)
	return expensiveItem, nil
}
