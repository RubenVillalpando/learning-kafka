package model

import (
	"encoding/json"
)

type Customer struct {
	ID    string
	Name  string
	Email string
	Age   int
}

func (c *Customer) Key() []byte {
	return []byte(c.ID)
}

func (c *Customer) Value() []byte {

	jsonBytes, err := json.Marshal(&c)

	if err != nil {
		return []byte{}
	}

	return jsonBytes
}
