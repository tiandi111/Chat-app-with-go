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

// Read messages from client
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		// push the message to message queue c.room.forward of room
		c.room.forward <- msg
	}
}

// Write messages to client
func (c *client) write() {
	defer c.socket.Close()
	// write all messages read from channel c.send to cilent
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
