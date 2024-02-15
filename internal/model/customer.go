package model

import (
	"encoding/json"
	"strconv"
)

type Customer struct {
	ID    int
	Name  string
	Email string
}

type CustomerValue struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *Customer) Key() []byte {
	return []byte(strconv.Itoa(c.ID))
}

func (c *Customer) Value() ([]byte, error) {
	val := &CustomerValue{
		Name:  c.Name,
		Email: c.Email,
	}

	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil

}
