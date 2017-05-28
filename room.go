package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan clientMsg
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan clientMsg),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run(closed func()) {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			log.Println("debug: New client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			log.Println("debug: Client left")
			if len(r.clients) == 0 {
				closed()
				return
			}
		case msg := <-r.forward:
			log.Println("debug: Message received: ", string(msg.message))
			for client := range r.clients {
				if client != msg.client {
					client.send <- msg.message
					log.Println("debug:  -- sent to client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
