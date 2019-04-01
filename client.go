package main 

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket	*websocket.Conn
	// This channel is used to send message from room to client
	send 	chan	[]byte 
	room 	*room
}

// Room read message from client
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// Room write message to client
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
