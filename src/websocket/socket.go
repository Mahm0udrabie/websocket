package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var hub = make(map[string]map[string]*websocket.Conn)
var MessageChan = make(chan Message)

// channel  (1,2,3) => [con1,con2]
func MessagesConn() http.Handler {
	return http.HandlerFunc(MessageSocket)
}
func MessageSocket(res http.ResponseWriter, req *http.Request) {
	//check if user send channel
	channel := req.FormValue("channel")
	if channel == "" {
		return
	}
	//check if channel is in our database
	//db conn, api => user
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true },
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}

	defer (func() {
		if err := recover(); err != nil {
			return
		}
	})()
	var conn, err = upgrader.Upgrade(res, req, nil)
	if err != nil {
		defer conn.Close()
		return
	}
	uuid, _ := uuid.NewRandom()
	connId := uuid.String()
	addConnToHub(channel, connId, conn)
	go readData(channel, connId, conn)
	for {
		select {
		case message := <-MessageChan:
			for _, connWrite := range hub[message.ToClient] {
				err := connWrite.WriteJSON(message)
				if err != nil {
					removeConnFromHub(channel, connId)
					defer connWrite.Close()
				}
			}
		}
	}
}
func addConnToHub(channel string, connId string, conn *websocket.Conn) {
	_, ok := hub[channel]
	if ok {
		hub[channel][connId] = conn
	} else {
		connHub := make(map[string]*websocket.Conn)
		hub[channel] = connHub
		hub[channel][connId] = conn
	}
}
func readData(channel string, connId string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			removeConnFromHub(channel, connId)
			defer conn.Close()
			return
		}
	}
}
func removeConnFromHub(channel string, connId string) {
	_, ok := hub[channel][connId]
	if ok {
		delete(hub[channel], connId)
	}
}

type Message struct {
	ToClient   string `json:"to_client"`
	FromClient string `json:"from_client"`
	Data       string `json:"data"`
	Type       string `json:"type"`
}
