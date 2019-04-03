package main

import (
	"github.com/chat/trace"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type room struct {
	// a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// tracer will receive trace information of activity
	// in the room
	tracer trace.Tracer
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func newRoom() *room {
	return &room{
		forward:	make(chan []byte),
		join:		make(chan *client),
		leave: 		make(chan *client),
		clients: 	make(map[*client]bool),
	}
}

func (r *room) run() {
	// In the for loop below, only one block will be ran 
	// at a time so that client map is only ever modified
	// by one thing at a time
	for {
		select {
			case client := <-r.join:
				r.clients[client] = true
				r.tracer.Trace("New client joined")
			case client := <-r.leave:
				delete(r.clients, client)
				close(client.send)
				r.tracer.Trace("Client left")
			case msg := <-r.forward:
				r.tracer.Trace("Message receivd: ", string(msg))
				for client := range r.clients {
					client.send <- msg
					r.tracer.Trace(" -- sent to client")
				}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client {
		socket: socket,
		send:  	make(chan []byte, messageBufferSize),
		room:	r,
	}
	r.join <- client
	defer func() { r.leave <- client } ()
	// Start a goroutine to write 
	go client.write()
	// Read from the current thread
	client.read()
}











