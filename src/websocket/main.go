package main

import (
	"net/http"
)

func main() {
	http.Handle("/messages/conn", MessagesConn())
	http.HandleFunc("/messages/send", SendMsg)
	http.ListenAndServe(":3000", nil)
}
