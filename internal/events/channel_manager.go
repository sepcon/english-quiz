package events

import (
	"fmt"
	"github.com/sepcon/quizprob/internal/common_errors"
	"github.com/sepcon/quizprob/pkg/model/event_service"
	"sync"

	"github.com/gorilla/websocket"
)

type ChannelID = string
type SubsriberID = string
type channelManagerImpl struct {
	subscribers           map[ChannelID]map[*websocket.Conn]SubsriberID
	userChannels          map[*websocket.Conn][]ChannelID
	subscriberConnections map[SubsriberID]*websocket.Conn
	mu                    sync.RWMutex
}

func NewChannelManager() *channelManagerImpl {
	return &channelManagerImpl{
		subscribers:           make(map[ChannelID]map[*websocket.Conn]SubsriberID),
		userChannels:          make(map[*websocket.Conn][]ChannelID),
		subscriberConnections: make(map[SubsriberID]*websocket.Conn),
	}
}

func (bs *channelManagerImpl) Connect(conn *websocket.Conn, subscriberID string) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if existingConn, ok := bs.subscriberConnections[subscriberID]; ok {
		conn = existingConn
	} else {
		bs.subscriberConnections[subscriberID] = conn
	}

}
func (bs *channelManagerImpl) Subscribe(subscriberID, channel string) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	var conn *websocket.Conn
	if existingConn, ok := bs.subscriberConnections[subscriberID]; ok {
		conn = existingConn
	} else {
		return common_errors.NewUserNotFoundError(subscriberID)
	}

	if _, ok := bs.subscribers[channel]; !ok {
		bs.subscribers[channel] = make(map[*websocket.Conn]SubsriberID)
	}

	bs.subscribers[channel][conn] = subscriberID
	bs.userChannels[conn] = append(bs.userChannels[conn], channel)
	return nil
}

func (bs *channelManagerImpl) Unsubscribe(channel, subscriberID string) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.unsubscribeUnsafe(channel, subscriberID)
}

func (bs *channelManagerImpl) CleanupObsoleteConnection(conn *websocket.Conn, subscriberID SubsriberID) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.onSubscriberConnectionClosedUnsafe(conn, subscriberID)
}

func (bs *channelManagerImpl) onSubscriberConnectionClosedUnsafe(conn *websocket.Conn, subscriberID SubsriberID) {
	delete(bs.subscriberConnections, subscriberID)
	for _, channel := range bs.userChannels[conn] {
		delete(bs.subscribers[channel], conn)
	}
	delete(bs.userChannels, conn)
}

func (bs *channelManagerImpl) unsubscribeUnsafe(channel, subscriberID string) {
	conn := bs.subscriberConnections[subscriberID]
	if _, ok := bs.subscribers[channel]; ok {
		delete(bs.subscribers[channel], conn)
	}

	for i, ch := range bs.userChannels[conn] {
		if ch == channel {
			bs.userChannels[conn] = append(bs.userChannels[conn][:i], bs.userChannels[conn][i+1:]...)
			break
		}
	}

	// If user has no more subscriptions, remove their connection
	if len(bs.userChannels[conn]) == 0 {
		delete(bs.subscriberConnections, subscriberID)
	}
}

func (bs *channelManagerImpl) Publish(channel event_service.ChannelIDType,
	msgBody event_service.MessageBodyType) error {
	msg := event_service.Message{
		Channel: channel,
		Body:    msgBody,
	}
	rawMsg, err := msg.Serialize()
	if err != nil {
		return fmt.Errorf("Cannot serialize data")
	}

	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if subscribers, ok := bs.subscribers[channel]; ok {
		for conn, subscriberID := range subscribers {
			err := conn.WriteMessage(websocket.BinaryMessage, rawMsg)
			if err != nil {
				fmt.Printf("\nError updating channel[%s] to subscriber[%s]:%v",
					channel, subscriberID, err)
				bs.CleanupObsoleteConnection(conn, subscriberID)
				conn.Close()
			}
		}
		return nil
	} else {
		return fmt.Errorf("Channel [%s] not found", channel)
	}

}
