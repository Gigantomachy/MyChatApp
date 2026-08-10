package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

const eventChannel = "chat.events"

type BrokerMessage struct {
	Users   []string        `json:"users"`
	Payload json.RawMessage `json:"payload"`
}

// broker fans WS events out across backend pods via Redis pub/sub - Cassandra remains the source of
// truth - clients refetch dropped events on reconnect
type Broker struct {
	client *redis.Client
}

func NewBroker(addr string) *Broker {
	return &Broker{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (b *Broker) Ping(ctx context.Context) error {
	if err := b.client.Ping(ctx).Err(); err != nil {
		log.Printf("redis ping failed: %v", err)
		return err
	}
	return nil
}

// sends an event envelope to the event channel
func (b *Broker) Publish(ctx context.Context, users []string, payload []byte) error {
	envelope, err := json.Marshal(BrokerMessage{Users: users, Payload: payload})
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, eventChannel, envelope).Err()
}

// go-redis handles reconnects and re-subscription automatically.
func (b *Broker) Subscribe(ctx context.Context, hub *Hub) {
	go func() {
		pubsub := b.client.Subscribe(ctx, eventChannel)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for msg := range ch {
			var envelope BrokerMessage
			if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
				log.Printf("dropping malformed ws event: %v", err)
				continue
			}
			hub.deliver(envelope.Users, envelope.Payload)
		}
	}()
}

func (b *Broker) Close() error {
	return b.client.Close()
}
