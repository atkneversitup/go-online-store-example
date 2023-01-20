package controller

import (
	"log"
	"myapp-go-echo/database"

	"github.com/labstack/echo"
)

type ProductController struct{}

func (pc *ProductController) GetProducts(c echo.Context) error {
	var products []database.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, products)
}
func (pc *ProductController) GetProductsById(c echo.Context) error {
	var product database.Product
	if err := database.DB.Find(&product).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, product)
}
func (pc *ProductController) CreateProducts(c echo.Context) error {
	var product database.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Create(&product).Error; err != nil {
		log.Println(err)
		return c.JSON(500, err)
	}
	return c.JSON(200, product)
}
func (pc *ProductController) UpdateProducts(c echo.Context) error {
	var product database.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Save(&product).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, product)
}
func (pc *ProductController) DeleteProducts(c echo.Context) error {
	var product database.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Delete(&product).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, product)
}
