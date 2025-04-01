package svc

import (
	"fmt"
	"frappuccino/helper"
	"frappuccino/internal/models"
)

func (s *svc) CreateInventory(data models.InventoryItem) error {
	if err := helper.IsValidName(data.Name); err != nil {
		s.Log.Error(err.Error(), "Invalid inventory name")
		return err
	}

	if err := helper.CheckerForInventItems(data); err != nil {
		s.Log.Error(err.Error(), "Invalid inventory data")
		return err
	}

	// Проверяем по Name и Unit
	existingItem, err := s.Repo.InventoryRepo.GetByNameAndUnit(data.Name, data.Unit)
	if err == nil {
		existingItem.Stock += data.Stock
		if err := s.Repo.InventoryRepo.PutInventory(existingItem.ID, existingItem); err != nil {
			s.Log.Error(err.Error(), "Failed to update existing inventory item")
			return err
		}
		s.Log.Info("Inventory item quantity updated", "id", existingItem.ID, "new_stock", existingItem.Stock)
		return nil
	} else if err.Error() != fmt.Sprintf("item with name %s and unit %s not found", data.Name, data.Unit) {
		s.Log.Error(err.Error(), "Error checking inventory item existence")
		return err
	}

	if err := s.Repo.InventoryRepo.CreateInventory(data); err != nil {
		s.Log.Error(err.Error(), "Failed to create inventory item")
		return err
	}

	s.Log.Info("Inventory item created successfully", "id", data.ID)
	return nil
}

func (s *svc) InventoriesGet() ([]models.InventoryItem, error) {
	data, err := s.Repo.InventoryRepo.GetInventory()
	if err != nil {
		s.Log.Error(err.Error(), "Failed to get inventories")
		return nil, err
	}

	s.Log.Info("Successfully fetched inventories", "count", len(data))
	return data, nil
}

func (s *svc) InventoryGetId(id int) (models.InventoryItem, error) {
	data, err := s.Repo.InventoryRepo.GetInventoryId(id)
	if err != nil {
		s.Log.Error(err.Error(), "Failed to fetch inventory item by ID")
		return models.InventoryItem{}, err
	}

	s.Log.Info("Successfully fetched inventory item", "id", id)
	return data, nil
}

func (s *svc) InventoryUpDate(id int, data models.InventoryItem) error {
	if err := helper.IsValidName(data.Name); err != nil {
		s.Log.Error(err.Error(), "Invalid inventory name")
		return err
	}

	if err := helper.CheckerForInventItems(data); err != nil {
		s.Log.Error(err.Error(), "Invalid inventory data")
		return err
	}

	if err := s.Repo.InventoryRepo.PutInventory(id, data); err != nil {
		s.Log.Error(err.Error(), "Failed to update inventory item")
		return err
	}

	s.Log.Info("Inventory item updated successfully", "id", id)
	return nil
}

func (s *svc) DeleteInvent(id int) error {
	if err := s.Repo.InventoryRepo.DeleteInvent(id); err != nil {
		s.Log.Error(err.Error(), "Failed to delete inventory item")
		return err
	}

	s.Log.Info("Inventory item deleted successfully", "id", id)
	return nil
}

func (s *svc) GetLeftOvers(sortBy string, page int, pageSize int) (*models.InventoryResponse, error) {
	// Получаем остатки инвентаря с репозитория
	items, err := s.Repo.InventoryRepo.GetLeftOversWithPagination(sortBy, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory leftovers: %w", err)
	}

	// Получаем общее количество элементов в инвентаре
	totalItems, err := s.Repo.InventoryRepo.CountTotalInventoryItems()
	if err != nil {
		return nil, fmt.Errorf("failed to count inventory items: %w", err)
	}

	// Рассчитываем количество страниц
	totalPages := (totalItems + pageSize - 1) / pageSize

	// Возвращаем данные с пагинацией
	return &models.InventoryResponse{
		CurrentPage: page,
		PageSize:    pageSize,
		HasNextPage: page*pageSize < totalItems,
		TotalPages:  totalPages,
		Data:        items, // Передаем срез []models.InventoryItem
	}, nil
}
