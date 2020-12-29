package model

import "time"

type Customer struct {
	Id           int       `json:"cID,omitempty"`
	Name         string    `json:"cName,omitempty"`
	Tel          int64     `json:"cTel,number,omitempty"`
	Address      string    `json:"cAddress,omitempty"`
	RegisterDate time.Time `json:"cRegisterDate,omitempty"`
}
