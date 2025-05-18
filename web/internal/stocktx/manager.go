package stocktx

import (
	"context"
	"fmt"
	"log"
)

type Manager struct {
	client Catalog
}

func NewManager(client Catalog) *Manager {
	return &Manager{client: client}
}

func (m *Manager) TryReserve(ctx context.Context, items []Item) ([]Item, error) {
	var reserved []Item

	for _, item := range items {
		err := m.client.DecreaseProductQuantity(ctx, item.ProductID, item.Amount)
		if err != nil {
			log.Printf("[stocktx] failed to reserve product %d (amount: %d): %v", item.ProductID, item.Amount, err)
			m.Rollback(ctx, reserved)
			return nil, fmt.Errorf("failed to reserve product %d: %w", item.ProductID, err)
		}
		reserved = append(reserved, item)
	}

	return reserved, nil
}

func (m *Manager) Rollback(ctx context.Context, items []Item) {
	for _, item := range items {
		err := m.client.IncreaseProductQuantity(ctx, item.ProductID, item.Amount)
		if err != nil {
			log.Printf("[stocktx] failed to rollback product %d (amount: %d): %v", item.ProductID, item.Amount, err)
		} else {
			log.Printf("[stocktx] rollback successful for product %d (amount: %d)", item.ProductID, item.Amount)
		}
	}
}
