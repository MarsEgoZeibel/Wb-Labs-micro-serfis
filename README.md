# WB Order Service
![wb](https://github.com/user-attachments/assets/f5e18f16-2b36-4ff8-a3ef-fe2958c59290)


Сервис для отображения данных о заказах с использованием NATS Streaming и PostgreSQL.

## Требования

- Go 1.21 или выше
- PostgreSQL 13 или выше
- NATS Streaming Server

## Установка и настройка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd wb-lab
```

2. Установите зависимости:
```bash
go mod download
```

3. Настройте PostgreSQL:
- Создайте базу данных:
```bash
createdb wb_orders
```
- Примените миграции:
```bash
psql -d wb_orders -f migrations/001_initial.sql
```

4. Запустите NATS Streaming Server:
```bash
docker run -d -p 4222:4222 -p 8222:8222 nats-streaming -cid test-cluster
```

5. Настройте конфигурацию в файле `.env`

6. Запустите сервис:
```bash
go run cmd/app/main.go
```

7. Для тестирования, запустите скрипт публикации:
```bash
go run scripts/publisher.go
```

## Структура проекта

```
wb-lab/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── cache/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── nats/
│   ├── repository/
│   └── service/
├── migrations/
├── scripts/
├── web/
│   ├── static/
│   │   ├── css/
│   │   └── js/
│   └── templates/
├── .env
└── go.mod
```

## API Endpoints

- `GET /` - Главная страница со списком заказов
- `GET /order/:id` - Страница с подробной информацией о заказе
- `GET /api/orders` - API endpoint для получения списка всех заказов
- `GET /api/orders/:id` - API endpoint для получения информации о конкретном заказе

## Функциональность

- Получение и сохранение заказов через NATS Streaming
- Сохранение данных в PostgreSQL
- In-memory кэширование с восстановлением из БД
- Веб-интерфейс для просмотра заказов
- REST API для получения данных
- Поиск заказов по ID
- Отображение детальной информации о заказе
- Просмотр JSON-представления заказа
