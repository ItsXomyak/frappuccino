package svc

import (
	"strings"

	"frappuccino/helper"
	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
)

func (s *svc) CreateMenuItem(item models.MenuItem) (*models.MenuItem, error) {
	s.Log.Info("Creating new menu item", "name", item.Name)
	if err := helper.IsValidName(item.Name); err != nil {
		s.Log.Error("Invalid menu item name", "error", err.Error())
		return nil, err
	}

	dataInvent, err := s.Repo.InventoryRepo.GetInventory()
	if err != nil {
		s.Log.Error("Failed to get existing inventory items", "error", err.Error())
		return nil, err
	}
	s.Log.Info("Inventory loaded", "count", len(dataInvent), "items", dataInvent) // Логируем инвентарь

	if len(item.Ingredients) > 0 {
		if err := helper.CheckerForMenuItems(item, dataInvent, item.Ingredients); err != nil {
			s.Log.Error("Invalid menu item ingredients", "error", err.Error())
			return nil, err
		}
	}

	createdItem, err := s.Repo.MenuRepo.CreateMenuItem(item)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			s.Log.Error("Menu item already exists", "name", item.Name)
			return nil, cerrors.ErrExist
		}
		s.Log.Error("Failed to create menu item in repository", "error", err.Error())
		return nil, err
	}

	for _, ing := range item.Ingredients {
		err = s.Repo.MenuRepo.AddIngredientToMenuItem(createdItem.ID, ing)
		if err != nil {
			s.Log.Error("Failed to add ingredient", "menu_item_id", createdItem.ID, "error", err.Error())
			return nil, err
		}
	}

	s.Log.Info("Successfully created menu item", "id", createdItem.ID)
	return createdItem, nil
}

func (s *svc) GetAllMenuItems() ([]models.MenuItem, error) {
	s.Log.Info("Retrieving all menu items")

	items, err := s.Repo.MenuRepo.GetAllMenuItems()
	if err != nil {
		s.Log.Error("Failed to retrieve menu items", "error", err.Error())
		return nil, err
	}

	s.Log.Info("Successfully retrieved menu items", "count", len(items))
	return items, nil
}

func (s *svc) GetMenuItemByID(id int) (*models.MenuItem, error) {
	s.Log.Info("Retrieving menu item by ID", "id", id)

	item, err := s.Repo.MenuRepo.GetMenuItemByID(id)
	if err != nil {
		s.Log.Error("Failed to retrieve menu item", "id", id, "error", err.Error())
		return nil, err // Репозиторий уже возвращает cerrors.ErrMenuItemNotFound
	}

	s.Log.Info("Successfully retrieved menu item", "id", id)
	return item, nil
}

func (s *svc) UpdateMenuItem(id int, item models.MenuItem) (*models.MenuItem, error) {
	s.Log.Info("Updating menu item", "id", id)

	if err := helper.IsValidName(item.Name); err != nil {
		s.Log.Error("Invalid menu item name", "error", err.Error())
		return nil, err
	}

	dataInvent, err := s.Repo.InventoryRepo.GetInventory()
	if err != nil {
		s.Log.Error("Failed to get existing inventory items", "error", err.Error())
		return nil, err
	}

	ingredients, err := s.Repo.MenuRepo.GetIngredientsByMenuItemID(item.ID)
	if err != nil {
		s.Log.Error("Failed to get ingredients", "menu_item_id", item.ID, "error", err.Error())
		return nil, err
	}

	if err := helper.CheckerForMenuItems(item, dataInvent, ingredients); err != nil {
		s.Log.Error("Invalid menu item ingredients", "error", err.Error())
		return nil, err
	}

	updatedItem, err := s.Repo.MenuRepo.UpdateMenuItem(id, item)
	if err != nil {
		s.Log.Error("Failed to update menu item", "id", id, "error", err.Error())
		return nil, err
	}

	s.Log.Info("Successfully updated menu item", "id", id)
	return updatedItem, nil
}

func (s *svc) DeleteMenuItem(id int) error {
	s.Log.Info("Deleting menu item", "id", id)

	err := s.Repo.MenuRepo.DeleteMenuItem(id)
	if err != nil {
		s.Log.Error("Failed to delete menu item", "id", id, "error", err.Error())
		return err // Репозиторий возвращает cerrors.ErrMenuItemNotFound
	}

	s.Log.Info("Successfully deleted menu item", "id", id)
	return nil
}
