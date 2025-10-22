package cache

import (
	"sync"
	"wb-lab/internal/models"
)

type Cache interface {
	Set(order models.Order)
	Get(orderUID string) (models.Order, bool)
	GetAll() []models.Order
}

type OrderCache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]models.Order),
	}
}

func (c *OrderCache) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *OrderCache) Get(orderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}

func (c *OrderCache) GetAll() []models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	orders := make([]models.Order, 0, len(c.orders))
	for _, order := range c.orders {
		orders = append(orders, order)
	}
	return orders
}
