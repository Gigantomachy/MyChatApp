package ws

import (
	"slices"
	"sync"
)

type Hub struct {
	mtx         sync.RWMutex
	connections map[string][]*Connection
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string][]*Connection),
	}
}

func (h *Hub) Register(c *Connection) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	h.connections[c.userID] = append(h.connections[c.userID], c)
}

func (h *Hub) Unregister(c *Connection) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if conns, ok := h.connections[c.userID]; ok {
		for i, connection := range conns {
			if connection.conn == c.conn {
				conns = slices.Delete(conns, i, i+1)
				h.connections[c.userID] = conns
				close(connection.send)
				if len(conns) == 0 {
					delete(h.connections, c.userID)
				}
				break
			}
		}
	}
}

func (h *Hub) PublishToUsers(user_ids []string, payload []byte) {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	for _, uid := range user_ids {
		for _, c := range h.connections[uid] {
			select {
			case c.send <- payload:
			default:
				// default case is here so that if the user has a full buffer we move on immediately after trying
				// we dont want to wait while holding lock
			}
		}
	}
}
