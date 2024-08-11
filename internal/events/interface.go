package events

import "github.com/gorilla/websocket"

type ChannelManager interface {
	Connect(conn *websocket.Conn, subscriberID string)
	Subscribe(subscriberID, channel string) error
	Unsubscribe(channel, subscriberID string)
	Publish(channel string, eventData []byte) (err error)
	CleanupObsoleteConnection(conn *websocket.Conn, subscriberID SubsriberID)
}
