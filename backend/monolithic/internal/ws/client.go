package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	pingWait  = (pongWait * 9) / 10
)

type Connection struct {
	hub    *Hub
	userID string
	conn   *websocket.Conn
	send   chan []byte
}

func (c *Connection) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	// conn.SetReadLimit ?
	// conn.SetReadDeadline ?
	// conn.SetPongHandler ?
	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break // clients send messages via REST API
		}
	}
}

func (c *Connection) writePump() {
	ticker := time.NewTicker(pingWait)
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
