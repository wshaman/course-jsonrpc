package methods

import (
	"encoding/json"
	"errors"
	"fmt"
)

type doHelloParams struct {
	Name string `json:"name"`
}

func DoHello(params json.RawMessage) (interface{}, error) {
	p := doHelloParams{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", errors.New("wrong params sent, expected {name:string}")
	}
	return fmt.Sprintf("Hello, %s! Gald to see you", p.Name), nil
}
