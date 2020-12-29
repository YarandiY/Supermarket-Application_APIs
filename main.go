package main

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/yarandiy/IE-assignment/handler"
	"github.com/yarandiy/IE-assignment/repository"
)

func init() {
	tmpDB, err := sql.Open("postgres", "dbname=my_database user=postgres password=admin host=127.0.0.1 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	repository.DB = tmpDB
}

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
