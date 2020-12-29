package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/yarandiy/IE-assignment/model"
)

var DB *sql.DB

func GetCustomer(input int) (model.Customer, error) {
	res := model.Customer{}

	var id int
	var name string
	var tel int64
	var address string
	var registerDate pq.NullTime

	err := DB.QueryRow(`SELECT id, name, tel, address, registerDate FROM customers where id = $1`, input).Scan(&id, &name, &tel, &address, &registerDate)
	if err == nil {
		res = model.Customer{Id: id, Name: name, Tel: tel, Address: address, RegisterDate: registerDate.Time}
	}

	return res, err
}

func AllCustomers() ([]model.Customer, error) {
	customers := []model.Customer{}

	rows, err := DB.Query(`SELECT id, name, tel, address, registerDate FROM customers order by id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var tel int64
		var address string
		var registerDate pq.NullTime

		err = rows.Scan(&id, &name, &tel, &address, &registerDate)
		if err != nil {
			return customers, err
		}

		currentCustomer := model.Customer{Id: id, Name: name, Tel: tel, Address: address}
		if registerDate.Valid {
			currentCustomer.RegisterDate = registerDate.Time
		}

		customers = append(customers, currentCustomer)
	}

	return customers, err
}

func InsertCustomer(name string, tel int64, address string, registerDate time.Time) (int, error) {
	var id int
	err := DB.QueryRow(`INSERT INTO customers(name, tel, address, registerDate) VALUES($1, $2, $3, $4) RETURNING id`, name, tel, address, registerDate).Scan(&id)

	if err != nil {
		return 0, err
	}

	fmt.Printf("Last inserted ID: %v\n", id)
	return id, err
}

func UpdateCustomer(id int, name string, tel int64, address string, registerDate time.Time) (int, error) {
	res, err := DB.Exec(`UPDATE customers set name=$1, tel=$2, address=$3, registerDate=$4 where id=$5 RETURNING id`, name, tel, address, registerDate, id)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsUpdated), err
}

func RemoveCustomer(id int) (int, error) {
	res, err := DB.Exec(`delete from customers where id = $1`, id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsDeleted), nil
}
