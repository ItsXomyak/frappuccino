package helper

import (
	"fmt"

	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
)

func CheckerForMenuItems(item models.MenuItem, allInventory []models.InventoryItem, ingredients []models.MenuItemIngredient) error {
	if len(ingredients) == 0 {
		return nil
	}
	inventoryMap := make(map[int]models.InventoryItem)
	for _, inv := range allInventory {
		inventoryMap[inv.ID] = inv
	}
	for _, i := range ingredients {
		if i.IngredientID == 0 {
			return fmt.Errorf("ingredient ID should not be empty")
		}
		if _, exists := inventoryMap[i.IngredientID]; !exists {
			return fmt.Errorf("ingredient with ID %d not found in inventory", i.IngredientID)
		}
		if i.Quantity <= 0 {
			return fmt.Errorf("ingredient quantity should be greater than 0, got: %f", i.Quantity)
		}
	}
	return nil
}

func CheckMenuExistsId(files []models.MenuItem, id int) error {
	if files == nil {
		return fmt.Errorf("menu items list is nil")
	}
	for _, i := range files {
		if i.ID == id {
			return nil
		}
	}
	return cerrors.ErrNotExist
}
