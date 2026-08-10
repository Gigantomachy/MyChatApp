package ws

import (
	"context"
	"log"
	"slices"
	"sync"
	"time"
)

type Hub struct {
	mtx         sync.RWMutex
	connections map[string][]*Connection
	broker      *Broker
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string][]*Connection),
	}
}

func (h *Hub) AttachBroker(b *Broker) {
	h.broker = b
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
	if h.broker == nil {
		h.deliver(user_ids, payload)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := h.broker.Publish(ctx, user_ids, payload); err != nil {
		// non-fatal: the message is already persisted in Cassandra, and clients refetch history on reconnect.
		log.Printf("failed to publish ws event to redis: %v", err)
	}
}

// deliver sends payload to the given users' connections on this process.
// it is called by PublishToUsers (local mode) and by the broker's Subscribe() in distributed mode
func (h *Hub) deliver(user_ids []string, payload []byte) {
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
