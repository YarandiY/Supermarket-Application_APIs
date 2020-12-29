package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/yarandiy/IE-assignment/handler"
)

func main() {
	e := echo.New()

	e.GET("/customers", handler.Customer{}.Read)
	e.DELETE("/customers/:id", handler.Customer{}.Delete)
	e.PUT("/customers/:id", handler.Customer{}.Update)
	e.POST("/customers", handler.Customer{}.Create)

	e.GET("/report/:month", handler.Report)

	if err := e.Start("0.0.0.0:8080"); err != nil {
		fmt.Println(err)
	}

}
