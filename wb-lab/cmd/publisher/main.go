package main

import (
	"log"
	"time"

	stan "github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "wb-publisher"
	subject := "orders"

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Попытка подключения к NATS Streaming (cluster: %s, client: %s)", clusterID, clientID)

	// Проверяем подключение к NATS несколько раз
	var sc stan.Conn
	var err error
	for i := 0; i < 3; i++ {
		sc, err = stan.Connect(
			clusterID,
			clientID,
			stan.NatsURL("nats://localhost:4222"),
			stan.ConnectWait(2*time.Second),
		)
		if err == nil {
			break
		}
		log.Printf("Попытка %d: Ошибка подключения: %v", i+1, err)
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Fatalf("Не удалось подключиться к NATS после 3 попыток: %v", err)
	}
	defer sc.Close()

	log.Println("Успешно подключились к NATS Streaming серверу")

	// Тестовый заказ
	orderJSON := `{
        "order_uid": "b563feb7b2b84b6test",
        "track_number": "WBILMTESTTRACK",
        "entry": "WBIL",
        "delivery": {
            "name": "Test Testov",
            "phone": "+9720000000",
            "zip": "2639809",
            "city": "Kiryat Mozkin",
            "address": "Ploshad Mira 15",
            "region": "Kraiot",
            "email": "test@gmail.com"
        },
        "payment": {
            "transaction": "b563feb7b2b84b6test",
            "request_id": "",
            "currency": "USD",
            "provider": "wbpay",
            "amount": 1817,
            "payment_dt": 1637907727,
            "bank": "alpha",
            "delivery_cost": 1500,
            "goods_total": 317,
            "custom_fee": 0
        },
        "items": [
            {
                "chrt_id": 9934930,
                "track_number": "WBILMTESTTRACK",
                "price": 453,
                "rid": "ab4219087a764ae0btest",
                "name": "Mascaras",
                "sale": 30,
                "size": "0",
                "total_price": 317,
                "nm_id": 2389212,
                "brand": "Vivienne Sabo",
                "status": 202
            }
        ],
        "locale": "en",
        "internal_signature": "",
        "customer_id": "test",
        "delivery_service": "meest",
        "shardkey": "9",
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    }`

	// Публикация сообщения
	log.Printf("Отправляем сообщение в канал %s", subject)
	log.Printf("Размер сообщения: %d байт", len(orderJSON))

	// Пытаемся опубликовать сообщение несколько раз
	for i := 0; i < 3; i++ {
		err = sc.Publish(subject, []byte(orderJSON))
		if err == nil {
			log.Printf("Сообщение успешно опубликовано!")
			// Ждем немного, чтобы подписчик успел обработать сообщение
			time.Sleep(time.Second)
			return
		}
		log.Printf("Попытка %d: Ошибка публикации: %v", i+1, err)
		time.Sleep(time.Second)
	}

	log.Fatalf("Не удалось опубликовать сообщение после 3 попыток")
}
