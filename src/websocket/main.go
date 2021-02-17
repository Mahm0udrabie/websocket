package main

import (
	"net/http"
)

func main (){
	http.Handle("messages/conn",MessagesConn())
	http.ListenAndServe("3000", nil)
}
