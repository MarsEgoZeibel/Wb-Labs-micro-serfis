package nats

import (
	"encoding/json"
	"log"
	"wb-lab/internal/models"

	stan "github.com/nats-io/stan.go"
)

type OrderHandler interface {
	HandleOrder(order models.Order) error
}

type NATSService struct {
	conn    stan.Conn
	handler OrderHandler
}

func NewNATSService(clusterID, clientID string, handler OrderHandler) (*NATSService, error) {
	conn, err := stan.Connect(clusterID, clientID)
	if err != nil {
		return nil, err
	}

	return &NATSService{
		conn:    conn,
		handler: handler,
	}, nil
}

func (s *NATSService) Subscribe(subject string) error {
	log.Printf("Подписка на канал: %s", subject)
	_, err := s.conn.Subscribe(subject, func(msg *stan.Msg) {
		log.Printf("Получено сообщение из NATS размером %d байт", len(msg.Data))
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("Ошибка разбора JSON сообщения: %v", err)
			return
		}
		log.Printf("Сообщение успешно разобрано, OrderUID: %s", order.OrderUID)

		err = s.handler.HandleOrder(order)
		if err != nil {
			log.Printf("Ошибка обработки заказа: %v", err)
			return
		}
		log.Printf("Заказ %s успешно обработан", order.OrderUID)
	}, stan.DurableName("wb-service-durable"))

	return err
}

func (s *NATSService) Close() error {
	return s.conn.Close()
}
