package main

import (
	"myapp-go-echo/routes"
	"os"

	//
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
)

func main() {
	os.Setenv("JWT_SECRET", "mysecret")
	e := echo.New()
	routes.InitRoutes(e)
	e.Logger.Fatal(e.Start("localhost:8080"))

}
