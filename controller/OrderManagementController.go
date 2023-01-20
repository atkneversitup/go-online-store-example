package controller

import (
	"fmt"
	"myapp-go-echo/database"
	"myapp-go-echo/database/form"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type OrderManagementController struct {
	Order            database.Order
	CostCalculator   OrderCostCalculator
	ItemAdder        OrderItemAdder
	PaymentProcessor OrderPaymentProcessor
}
type OrderItemAdderImpl struct {
}

type OrderItemAdder interface {
	AddItem(orderProduct database.OrderProduct) error
	RemoveItem(item database.OrderProduct) error
}
type OrderCostCalculator interface {
	// CalculateTotalCost() float64
}
type OrderPaymentProcessor interface {
	// ProcessPayment(order Order) error
}

// itemAdder func
func (oi OrderItemAdderImpl) AddItem(orderProduct database.OrderProduct) error {
	// Check if the order product already exists
	existingOrderProduct := database.OrderProduct{}
	if err := database.DB.Where("order_id = ? AND product_id = ?", orderProduct.OrderID, orderProduct.ProductID).First(&existingOrderProduct).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the order product does not exist, create a new one
			if err := database.DB.Create(&orderProduct).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// If the order product already exists, update the quantity
		existingOrderProduct.Quantity += orderProduct.Quantity
		if err := database.DB.Save(&existingOrderProduct).Error; err != nil {
			return err
		}
	}
	return nil
}
func (oi OrderItemAdderImpl) RemoveItem(orderProduct database.OrderProduct) error {
	// Check if the order product already exists
	existingOrderProduct := database.OrderProduct{}
	if err := database.DB.Where("order_id = ? AND product_id = ?", orderProduct.OrderID, orderProduct.ProductID).First(&existingOrderProduct).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the order product does not exist, return an error
			return fmt.Errorf("Order product not found")
		} else {
			return err
		}
	} else {
		// If the order product already exists, delete the record
		if err := database.DB.Delete(&existingOrderProduct).Error; err != nil {
			return err
		}
	}
	return nil
}

type OrderResponse struct {
	Order         OrderDetailResponse    `json:"Order"`
	OrderProducts []OrderProductResponse `json:"OrderProducts"`
}
type OrderDetailResponse struct {
	ID       uint    `json:"ID"`
	UserId   int     `json:"UserID"`
	StatusId uint    `json:"StatusID"`
	Summary  float32 `json:"Summary"`
}
type OrderProductResponse struct {
	ID        uint `json:"ID"`
	ProductID uint `json:"ProductID"`
	Name      string
	Price     float32
	Quantity  int `json:"Quantity"`
}

// create order // and status_id fixed Draft and order date is current date
func (omc *OrderManagementController) CreateOrder(c echo.Context) error {
	var fd form.FormOrder
	var status database.Status
	if err := c.Bind(&fd); err != nil {
		return c.JSON(400, err)
	}
	err := database.DB.Find(&status, "Name = ?", fd.Status).Error
	if err != nil {
		return c.JSON(500, err)
	}
	order := database.Order{
		UserID:   int(fd.UserID),
		StatusID: uint(status.ID),
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, order)
}
func (omc *OrderManagementController) GetAllOrders(c echo.Context) error {
	var orders []database.Order
	if err := database.DB.Find(&orders).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, orders)
}

// GetOrderById
func (omc *OrderManagementController) GetOrderById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, err)
	}
	claims := c.Get("claims").(jwt.MapClaims)
	user_id := claims["id"].(float64)

	var order database.Order
	order.ID = uint(id)
	order.UserID = int(user_id)
	// jsonObject, _ := json.Marshal(order)
	// log.Println(string(jsonObject))

	if err := database.DB.First(&order).Related(&order.OrderProducts, "OrderProducts").Error; err != nil {
		return c.JSON(500, err.Error())
	}
	var summary float32 = 0
	var ops []database.OrderProduct
	if err := database.DB.Find(&ops, "order_id = ?", order.ID).Error; err != nil {
		return c.JSON(500, err)
	}
	var orderProductsResponse []OrderProductResponse
	for _, op := range ops {
		if err := database.DB.Model(op).Related(&op.Product).Error; err != nil {
			return c.JSON(500, err)
		}
		summary += op.Product.Price * float32(op.Quantity)
		orderProductsResponse = append(orderProductsResponse, OrderProductResponse{
			ID:        op.ID,
			ProductID: op.ProductID,
			Name:      op.Product.Name,
			Price:     op.Product.Price,
			Quantity:  op.Quantity,
		})
	}
	// start calculate summary
	OrderDetailResponse := OrderDetailResponse{
		ID:       order.ID,
		UserId:   order.UserID,
		StatusId: order.StatusID,
		Summary:  summary,
	}
	orderResponse := OrderResponse{
		Order:         OrderDetailResponse,
		OrderProducts: orderProductsResponse,
	}

	return c.JSON(200, orderResponse)
}

func (omc *OrderManagementController) UpdateOrderById(c echo.Context) error {
	var fd form.FormOrder
	var status database.Status
	var order database.Order
	if err := c.Bind(&fd); err != nil {
		return c.JSON(400, err)
	}
	err := database.DB.Find(&status, "Name = ?", fd.Status).Error
	if err != nil {
		return c.JSON(500, err)
	}
	order.StatusID = uint(status.ID)
	order.UserID = int(fd.UserID)
	id := c.Param("id")
	var count int
	if err := database.DB.Model(&order).Where("id = ?", id).Count(&count).Error; err != nil {
		return c.JSON(500, err)
	}
	if count == 0 {
		return c.JSON(404, "record not found")
	}
	if err := database.DB.Model(&order).Where("id = ?", id).Updates(map[string]interface{}{"status_id": order.StatusID, "user_id": order.UserID}).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, order)
}

func (omc *OrderManagementController) AddProductToOrder(c echo.Context) error {
	var order database.Order
	var product database.Product
	var orderProduct database.OrderProduct
	// Get the user ID and admin value from the JWT token
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// user_id := int(claims["id"].(float64))
	if err := c.Bind(&orderProduct); err != nil {
		return c.JSON(400, err)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, err)
	}
	var count int
	if err := database.DB.Model(&order).Where("id = ?", id).Count(&count).Error; err != nil {
		return c.JSON(500, err)
	}
	if count == 0 {
		return c.JSON(404, "record order not found")
	}
	omc.Order = order
	if err := database.DB.Model(&product).Where("id = ?", orderProduct.ProductID).Count(&count).Error; err != nil {
		return c.JSON(500, err)
	}
	if count == 0 {
		return c.JSON(404, "record product not found")
	}
	orderProduct.OrderID = uint(id)
	if err := omc.ItemAdder.AddItem(orderProduct); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(200, orderProduct)
}
