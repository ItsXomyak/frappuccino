package svc

import (
	"fmt"
	"frappuccino/internal/models"
	"strings"
)

func (s *svc) SearchFullText(query, filter, minPrice, maxPrice string) (*models.SearchResponse, error) {
	menuResults, err := s.Repo.SearchRepo.SearchMenuItems(query, filter, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to search menu items: %w", err)
	}

	orderResults, err := s.Repo.SearchRepo.SearchOrders(query, filter, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to search orders: %w", err)
	}

	// Рассчитываем релевантность для каждого результата ( надо норм сделать а не так)
	for i := range menuResults {
		menuResults[i].Relevance = calculateRelevance(menuResults[i].Name, query)
	}

	for i := range orderResults {
		orderResults[i].Relevance = calculateRelevance(orderResults[i].Name, query)
	}

	// Формируем ответ
	response := &models.SearchResponse{
		MenuItems:    menuResults,
		Orders:       orderResults,
		TotalMatches: len(menuResults) + len(orderResults),
	}

	return response, nil
}

func calculateRelevance(name, query string) float64 {
	if strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
		// если строка совпала, релевантность = 0.8
		return 0.8
	}
	return 0.0
}

func (s *svc) GetOrderedItemsByPeriod(period, month, year string) ([]models.OrderedItemReport, error) {
	report, err := s.Repo.SearchRepo.GetOrderedItemsByPeriod(period, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered items by period: %w", err)
	}

	return report, nil
}
