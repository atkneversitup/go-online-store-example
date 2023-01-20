package controller

import (
	"myapp-go-echo/database"

	"github.com/labstack/echo"
)

type OrderController struct{}

func (oc *OrderController) GetOrders(c echo.Context) error {
	var orders []database.Order
	if err := database.DB.Find(&orders).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, orders)
}
func (oc *OrderController) GetOrdersById(c echo.Context) error {
	var order database.Order
	if err := database.DB.Find(&order).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, order)
}
func (oc *OrderController) CreateOrders(c echo.Context) error {
	var order database.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, order)
}
func (oc *OrderController) UpdateOrders(c echo.Context) error {
	var order database.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Save(&order).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, order)
}
func (oc *OrderController) DeleteOrders(c echo.Context) error {
	var order database.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Delete(&order).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, order)
}
