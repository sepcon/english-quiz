package event_sdk

import (
	"context"
	"github.com/sepcon/quizprob/pkg/model/event_service"
)

type ConnectionErrorCallback = func(error)
type ProcessEventCallback = func([]byte)

type MessageHandler interface {
	OnConnectionError(err error)
	OnMessage([]byte)
}

type Publisher interface {
	Publish(ctx context.Context, channel event_service.ChannelIDType, msgBody event_service.MessageBodyType) error
}

type Subscriber interface {
	Connect(ctx context.Context) error
	Disconnect()
	Subscribe(ctx context.Context, channel event_service.ChannelIDType, handler MessageHandler) error
	Unsubscribe(ctx context.Context, channel event_service.ChannelIDType) error
}
