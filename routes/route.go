package routes

import (
	"myapp-go-echo/controller"
	"myapp-go-echo/middleware"

	"github.com/labstack/echo"
)

func InitRoutes(e *echo.Echo) {
	// User routes
	userController := &controller.UserController{}
	productController := &controller.ProductController{}
	// orderController := &controller.OrderController{}
	oia := &controller.OrderItemAdderImpl{}
	orderManagementController := &controller.OrderManagementController{
		ItemAdder: oia,
	}

	e.GET("/users", userController.GetUsers)
	e.POST("/users", userController.CreateUsers)
	e.POST("/register", userController.Register)
	e.POST("/login", userController.Login)

	e.GET("/products", productController.GetProducts)
	e.GET("/products/:id", productController.GetProductsById)
	e.POST("/products", productController.CreateProducts)
	e.PUT("/products/:id", productController.UpdateProducts)
	e.DELETE("/products/:id", productController.DeleteProducts)

	e.GET("/orders", middleware.JwtMiddleware(orderManagementController.GetAllOrders))
	e.GET("/orders/:id", middleware.JwtMiddleware(orderManagementController.GetOrderById))
	e.POST("/orders", middleware.JwtMiddleware(orderManagementController.CreateOrder))
	e.PATCH("/orders/:id", middleware.JwtMiddleware(orderManagementController.UpdateOrderById))
	e.PUT(("orders/:id/orderproduct"), middleware.JwtMiddleware(orderManagementController.AddProductToOrder))

}
