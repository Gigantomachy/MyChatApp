package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true, // frontend, we mirror the CORS config
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return allowedOrigins[r.Header.Get("Origin")]
	},
}

func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		connection := &Connection{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			userID: userID,
		}

		hub.Register(connection)

		go connection.writePump()
		go connection.readPump()
	}
}
