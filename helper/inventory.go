package helper

import (
	"fmt"
	"strings"

	"frappuccino/internal/models"
	"frappuccino/pkg/cerrors"
)

func CheckItemExists(items []models.InventoryItem, name, unit string) (*models.InventoryItem, error) {
	if items == nil {
		return nil, fmt.Errorf("inventory items list is nil")
	}

	for _, item := range items {
		if strings.ToLower(item.Name) == strings.ToLower(name) && strings.ToLower(item.Unit) == strings.ToLower(unit) {
			return &item, nil
		}
	}
	return nil, cerrors.ErrNotExist
}

func CheckItemId(items []models.InventoryItem, id int) error {
	if items == nil {
		return fmt.Errorf("inventory items list is nil")
	}

	for _, item := range items {
		if item.ID == id {
			return nil
		}
	}
	return fmt.Errorf("item with ID %d not found", id)
}

func CheckerForInventItems(item models.InventoryItem) error {
	if item.ID != 0 {
		return fmt.Errorf("item ID should not be set when adding a new item")
	}
	if item.Name == "" {
		return fmt.Errorf("please provide a name for the item")
	}
	if item.Stock <= 0 {
		return fmt.Errorf("please specify a stock quantity for the item")
	}
	if item.Unit == "" {
		return fmt.Errorf("please provide a unit for the item")
	}
	return nil
}
