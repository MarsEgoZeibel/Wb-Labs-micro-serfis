package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	sc, err := stan.Connect("test-cluster", "test-publisher", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	cities := []string{"Москва", "Санкт-Петербург", "Казань", "Новосибирск", "Екатеринбург"}
	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "Under Armour"}
	items := []string{"Кроссовки", "Футболка", "Шорты", "Куртка", "Носки"}

	for i := 1; i <= 10; i++ {
		order := Order{
			OrderUID:    fmt.Sprintf("test_order_%d", i),
			TrackNumber: fmt.Sprintf("TRACK%d", i),
			Entry:       "WBIL",
			Delivery: Delivery{
				Name:    fmt.Sprintf("Покупатель %d", i),
				Phone:   fmt.Sprintf("+7900555%04d", i),
				Zip:     fmt.Sprintf("1%04d0", i),
				City:    cities[i%len(cities)],
				Address: fmt.Sprintf("ул. Примерная, д. %d", i),
				Region:  "Московская область",
				Email:   fmt.Sprintf("customer%d@example.com", i),
			},
			Payment: Payment{
				Transaction:  fmt.Sprintf("tr_%d", i),
				RequestID:    "",
				Currency:     "RUB",
				Provider:     "wbpay",
				Amount:       1000 * i,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1000,
				GoodsTotal:   1000 * i,
				CustomFee:    0,
			},
			Items: []Item{
				{
					ChrtID:      i,
					TrackNumber: fmt.Sprintf("TRACK%d", i),
					Price:       1000 * i,
					RID:         fmt.Sprintf("rid_%d", i),
					Name:        items[i%len(items)],
					Sale:        0,
					Size:        "M",
					TotalPrice:  1000 * i,
					NmID:        i,
					Brand:       brands[i%len(brands)],
					Status:      202,
				},
			},
			Locale:            "ru",
			InternalSignature: "",
			CustomerID:        fmt.Sprintf("customer_%d", i),
			DeliveryService:   "meest",
			ShardKey:          "9",
			SmID:              99,
			DateCreated:       time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			OofShard:          "1",
		}

		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Ошибка маршалинга заказа %d: %v", i, err)
			continue
		}

		err = sc.Publish("orders", orderJSON)
		if err != nil {
			log.Printf("Ошибка публикации заказа %d: %v", i, err)
			continue
		}

		log.Printf("Опубликован заказ %d", i)
		time.Sleep(100 * time.Millisecond)
	}
}
