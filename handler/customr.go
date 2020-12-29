package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yarandiy/IE-assignment/model"
	// "go.mongodb.org/mongo-driver/mongo"
)

type Customer struct {
}

type Customers []model.Customer

var customers Customers

// var client *mongo.Client

type Request struct {
	Name    string `json:"cName"`
	Tel     int64  `json:"cTel,number"`
	Address string `json:"cAddress"`
}

type CreateResponse struct {
	Id           int       `json:"cID,omitempty"`
	Name         string    `json:"cName,omitempty"`
	Tel          int64     `json:"cTel,number,omitempty"`
	Address      string    `json:"cAddress,omitempty"`
	RegisterDate time.Time `json:"cRegisterDate,omitempty"`
	Message      string    `json:"msg,omitempty"`
}

func (customer Customer) Create(c echo.Context) error {

	var req Request

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.Name == "" || req.Address == "" || req.Tel == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, SimpleResponse{
			Message: "error : input validation error",
		})
	}

	m := model.Customer{
		Name:         req.Name,
		Tel:          req.Tel,
		Address:      req.Address,
		RegisterDate: time.Now(),
		Id:           rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000),
	}

	customers = append(customers, m)

	return c.JSON(http.StatusCreated, CreateResponse{
		Name:         m.Name,
		Tel:          m.Tel,
		Address:      m.Address,
		RegisterDate: m.RegisterDate,
		Id:           m.Id,
		Message:      "success",
	})
}

type ReadResponse struct {
	Size      int       `json:"size"`
	Customers Customers `json:"customers,omitempty"`
	Message   string    `json:"msg,omitempty"`
}

func (customer Customer) Read(c echo.Context) error {
	if len(c.QueryString()) != 0 {
		value := c.QueryParam("cName")

		if value == "" {
			return c.JSON(http.StatusBadRequest, SimpleResponse{
				Message: "error : query param validation error",
			})
		}

		for _, custo := range customers {
			if strings.Contains(custo.Name, value) {
				return c.JSON(http.StatusOK, custo)
			}
		}
		return c.JSON(http.StatusBadRequest, SimpleResponse{
			Message: "error : there is no customer with this cName",
		})
	}
	c_size := len(customers)
	msg := "success"
	if customers == nil || c_size == 0 {
		msg = "error : there is no customer"
	}
	return c.JSON(http.StatusOK, ReadResponse{
		Size:      c_size,
		Customers: customers,
		Message:   msg,
	})
}

type SimpleResponse struct {
	Message string `json:"msg,omitempty"`
}

func (customer Customer) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, _ := range customers {
		if customers[i].Id == id {
			customers = append(customers[:i], customers[i+1:]...)
			return c.JSON(http.StatusOK, SimpleResponse{
				Message: "success",
			})
		}
	}
	return echo.NewHTTPError(http.StatusBadRequest, SimpleResponse{
		Message: "error : the input cID is not available ",
	})
}

func (customer Customer) Update(c echo.Context) error {

	var req Request

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, SimpleResponse{
			Message: "error : input validation error",
		})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	for i, _ := range customers {
		if customers[i].Id == id {
			if req.Name != "" {
				customers[i].Name = req.Name
			}
			if req.Tel != 0 {
				customers[i].Tel = req.Tel
			}
			if req.Address != "" {
				customers[i].Address = req.Address
			}
			return c.JSON(http.StatusOK, CreateResponse{
				Name:         customers[i].Name,
				Tel:          customers[i].Tel,
				Address:      customers[i].Address,
				RegisterDate: customers[i].RegisterDate,
				Id:           customers[i].Id,
				Message:      "success",
			})
		}
	}
	return c.JSON(http.StatusBadRequest, SimpleResponse{
		Message: "error : the input cID is not available ",
	})
}

type ReportResponse struct {
	Total   int    `json:"totalCustomers"`
	Period  int    `json:"period"`
	Message string `json:"msg,omitempty"`
}

func Report(c echo.Context) error {
	month, _ := strconv.Atoi(c.Param("month"))
	counter := 0
	for _, custo := range customers {
		m := int(custo.RegisterDate.Month()) - 1
		if m == month {
			counter += 1
		}
	}
	if counter == 0 {
		return c.JSON(http.StatusBadRequest, SimpleResponse{
			Message: "error : there is no customer in this period",
		})
	}
	return c.JSON(http.StatusBadRequest, ReportResponse{
		Total:   counter,
		Period:  month,
		Message: "success",
	})
}
