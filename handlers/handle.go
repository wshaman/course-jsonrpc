package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wshaman/course-jsonrpc/handlers/methods"
)

type Request struct {
	Method string `json:"method"`
	Params json.RawMessage `json:"params"`
	Id int64 `json:"id"`
}

type Response struct {
	Id int64 `json:"id"`
	Error interface{} `json:"error"`
	Result interface{} `json:"result"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	bodyData, err := io.ReadAll(body)
	if err != nil {
		log.Print(err.Error())
		returnErr(w, -1, "Failed to read request body")
		return
	}
	req := &Request{}
	if err = json.Unmarshal(bodyData, req); err != nil {
		log.Print(err.Error())
		returnErr(w, -1, "Request format is wrong")
		return
	}
	handlers := map[string]func(message json.RawMessage)(interface{}, error){
		"doHello": methods.DoHello,
	}
	h, ok := handlers[req.Method]
	if !ok {
		returnErr(w, req.Id, fmt.Sprintf("method %s is not supported", req.Method))
		return
	}
	resp, err := h(req.Params)
	if err != nil {
		returnErr(w, req.Id, err.Error())
		return
	}
	returnOK(w, req.Id, resp)
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