package service

import (
	"log"
	"wb-lab/internal/cache"
	"wb-lab/internal/models"
	"wb-lab/internal/repository"
)

type OrderService struct {
	repo  *repository.OrderRepository
	cache cache.Cache
}

func NewOrderService(repo *repository.OrderRepository, cache cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) HandleOrder(order models.Order) error {
	log.Printf("Получен новый заказ: %s", order.OrderUID)

	err := s.repo.SaveOrder(order)
	if err != nil {
		log.Printf("Ошибка сохранения заказа в БД: %v", err)
		return err
	}
	log.Printf("Заказ %s сохранен в БД", order.OrderUID)

	// Сохраняем заказ в кэш
	s.cache.Set(order)
	log.Printf("Заказ %s сохранен в кэш", order.OrderUID)
	return nil
}

func (s *OrderService) GetOrder(orderUID string) (models.Order, bool) {
	return s.cache.Get(orderUID)
}

func (s *OrderService) GetAllOrders() []models.Order {
	return s.cache.GetAll()
}

func (s *OrderService) RestoreCache() error {
	orders, err := s.repo.GetAllOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.cache.Set(order)
	}

	return nil
}
