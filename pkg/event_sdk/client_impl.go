package event_sdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sepcon/quizprob/pkg/model/event_service"
	"net/http"
	"net/url"
	"sync"
)

func NewSubscriber(baseURL string, clientID string) Subscriber {
	return &subscriberImpl{
		client:   client{baseURL: baseURL},
		userID:   clientID,
		handlers: make(map[event_service.ChannelIDType]MessageHandler),
	}
}

func NewPublisher(baseURL string) Publisher {
	return &publissherImpl{
		client{baseURL: baseURL},
	}
}

type client struct {
	baseURL string
}

type subscriberImpl struct {
	client
	userID   string
	handlers map[event_service.ChannelIDType]MessageHandler
	conn     *websocket.Conn
	mu       sync.Mutex
}

type publissherImpl struct {
	client
}

func (bc *subscriberImpl) Disconnect() {
	if bc.conn != nil {
		bc.conn.Close()
	}
}

func (bc *subscriberImpl) Connect(ctx context.Context) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.conn != nil {
		return fmt.Errorf("already connected")
	}

	u := url.URL{Scheme: "ws", Host: bc.baseURL, Path: "/connect"}
	q := u.Query()
	q.Set("subscriber_id", bc.userID)
	u.RawQuery = q.Encode()

	var err error
	bc.conn, _, err = websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	// Start a goroutine to handle incoming messages
	go bc.handleMessages()

	return nil
}

func (bc *subscriberImpl) Subscribe(ctx context.Context, channel event_service.ChannelIDType, handler MessageHandler) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.handlers[channel] != nil {
		return fmt.Errorf("already subscribed for channel %s!", channel)
	}

	if bc.conn == nil {
		return fmt.Errorf("not connected")
	}

	err := bc.conn.WriteJSON(map[string]string{
		"type":    "subscribe",
		"channel": channel,
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}

	bc.handlers[channel] = handler
	return nil
}

func (bc *subscriberImpl) Unsubscribe(ctx context.Context, channel event_service.ChannelIDType) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.conn == nil {
		return fmt.Errorf("not connected")
	}

	err := bc.conn.WriteJSON(map[string]string{
		"type":    "unsubscribe",
		"channel": channel,
	})
	if err != nil {
		return fmt.Errorf("failed to unsubscribe: %v", err)
	}
	return nil
}

func (bc *subscriberImpl) handleMessages() {
	for {
		_, rawMessage, err := bc.conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		subscriberMsg := event_service.Message{}
		err = subscriberMsg.Deserialize(rawMessage)
		if err != nil {
			println("Error: cannot deserialized data", err.Error())
		}

		if handler := bc.handlers[subscriberMsg.Channel]; handler != nil {
			handler.OnMessage(subscriberMsg.Body)
		}
	}
}

func (bc *publissherImpl) Publish(ctx context.Context, channel event_service.ChannelIDType, msgBody event_service.MessageBodyType) error {
	u := url.URL{Scheme: "http", Host: bc.baseURL, Path: "/publish/" + url.QueryEscape(url.PathEscape(channel))}

	resp, err := http.Post(u.String(), "application/octet-stream", bytes.NewBuffer(msgBody))

	if err != nil {
		return fmt.Errorf("failed to publish event: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to publish event, status code: %d", resp.StatusCode)
	}

	return nil
}
