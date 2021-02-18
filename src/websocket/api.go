package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMsg(r http.ResponseWriter, q *http.Request) {
	var MessageRequest Message
	q.Body=http.MaxBytesReader(r,q.Body,1048576)
	dec:=json.NewDecoder(q.Body)
	err:= dec.Decode(&MessageRequest)
	if err != nil {
		fmt.Println("error decode body")
		return
	}
	MessageChan <- MessageRequest
	var Response struct {
		Status bool `json:"status"`
	}
	Response.Status=true
	message,err:=json.Marshal(Response)
	r.Write(message)
}
