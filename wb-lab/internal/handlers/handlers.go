package handlers

import (
	"log"
	"net/http"
	"wb-lab/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderUID := c.Param("id")
	order, exists := h.service.GetOrder(orderUID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders := h.service.GetAllOrders()
	log.Printf("Получено заказов: %d", len(orders))
	c.JSON(http.StatusOK, orders)
}

func SetupRoutes(router *gin.Engine, handler *OrderHandler) {
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/order/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", nil)
	})

	api := router.Group("/api")
	{
		api.GET("/orders", handler.GetAllOrders)
		api.GET("/orders/:id", handler.GetOrder)
	}
}
