package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yarandiy/IE-assignment/model"
	"github.com/yarandiy/IE-assignment/repository"
)

type Customer struct {
}

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

	t := time.Now()
	id, err := repository.InsertCustomer(req.Name, req.Tel, req.Address, t)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	m := model.Customer{
		Name:         req.Name,
		Tel:          req.Tel,
		Address:      req.Address,
		RegisterDate: t,
		Id:           id,
	}

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
	Size      int              `json:"size"`
	Customers []model.Customer `json:"customers,omitempty"`
	Message   string           `json:"msg,omitempty"`
}

func (customer Customer) Read(c echo.Context) error {
	if len(c.QueryString()) != 0 {
		value := c.QueryParam("cName")

		if value == "" {
			return c.JSON(http.StatusBadRequest, SimpleResponse{
				Message: "error : query param validation error",
			})
		}
		customers, err := repository.AllCustomers()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		counter := 0
		var result []model.Customer
		for _, custo := range customers {
			if strings.HasPrefix(custo.Name, value) {
				counter += 1
				result = append(result, custo)
			}
		}
		if counter != 0 {
			return c.JSON(http.StatusOK, ReadResponse{
				Size:      counter,
				Customers: result,
				Message:   "success",
			})
		}
		return c.JSON(http.StatusBadRequest, SimpleResponse{
			Message: "error : there is no customer with this cName",
		})
	}

	customers, err := repository.AllCustomers()
	c_size := len(customers)
	msg := "success"
	if customers == nil || c_size == 0 {
		msg = "error : there is no customer"
	}
	if err != nil {
		msg = "error : database error"
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
	rows, err := repository.RemoveCustomer(id)
	if err != nil || rows == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, SimpleResponse{
			Message: "error : the input cID is not available ",
		})
	}
	return c.JSON(http.StatusOK, SimpleResponse{
		Message: "success",
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
	custo, err := repository.GetCustomer(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, SimpleResponse{
			Message: "error : the input cID is not available ",
		})
	}
	name := req.Name
	tel := req.Tel
	address := req.Address

	if name == "" {
		name = custo.Name
	}
	if tel == 0 {
		tel = custo.Tel
	}
	if req.Address == "" {
		address = custo.Address
	}

	_, err = repository.UpdateCustomer(id, name, tel, address, custo.RegisterDate)

	if err != nil {
		return c.JSON(http.StatusBadRequest, SimpleResponse{
			Message: "error : the input cID is not available ",
		})
	}

	return c.JSON(http.StatusOK, CreateResponse{
		Name:         name,
		Tel:          tel,
		Address:      address,
		RegisterDate: custo.RegisterDate,
		Id:           id,
		Message:      "success",
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
	customers, err := repository.AllCustomers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
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
