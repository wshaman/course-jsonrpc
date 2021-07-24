package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wshaman/course-jsonrpc/handlers/methods"
)

type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Id     int64           `json:"id"`
}

type Response struct {
	Id     int64       `json:"id"`
	Error  interface{} `json:"error"`
	Result interface{} `json:"result"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		returnErr(w, -1, err.Error())
		return
	}
	resp, err := handle(*req)
	if err != nil {
		returnErr(w, req.Id, err.Error())
		return
	}
	returnOK(w, req.Id, resp)
}

func parseRequest(r *http.Request) (*Request, error) {
	body := r.Body
	defer body.Close()
	bodyData, err := io.ReadAll(body)
	if err != nil {
		log.Print(err.Error())
		return nil, errors.New("failed to read request body")
	}
	req := &Request{}
	if err = json.Unmarshal(bodyData, req); err != nil {
		log.Print(err.Error())
		return nil, errors.New("request format is wrong")
	}
	return req, nil
}

func handle(req Request) (interface{}, error) {
	handlers := map[string]func(message json.RawMessage) (interface{}, error){
		"doHello":    methods.DoHello,
		"summ":       methods.Summ,
		"summStatus": methods.SummStatus,
	}
	h, ok := handlers[req.Method]
	if !ok {
		return nil, errors.New(fmt.Sprintf("method %s is not supported", req.Method))
	}
	resp, err := h(req.Params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func returnOK(w http.ResponseWriter, id int64, data interface{}) {
	res := Response{
		Id:     id,
		Error:  nil,
		Result: data,
	}
	res.writeToWeb(w)
}

func returnErr(w http.ResponseWriter, id int64, data interface{}) {
	res := Response{
		Id:     id,
		Error:  data,
		Result: nil,
	}
	res.writeToWeb(w)
}

func (r Response) writeToWeb(w http.ResponseWriter) {
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}
