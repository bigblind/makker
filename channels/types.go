package channels

import (
	"context"
	"net/http"
)

type Channel interface {
	Public() bool
	Namespace() string
	Id() string
	ClientId() string
	Emit(event string, data interface{})
	EmitExcluding(socketId, event string, data interface{})
}

type Event struct {
	Channel Channel
	Name string
	Data interface{}
}

type EventHandler func(ctx context.Context, channel Channel, userId, socketId string)
type ChannelAuthChecker func(ctx context.Context, channel Channel, userId string) error

type ChannelProvider interface {
	NewChannel(ctx context.Context, namespace, id string, public bool) Channel
	EmitBatch(ctx context.Context, events []Event)
	OnJoin(namespace string, handler EventHandler)
	OnLeave(namespace string, handler EventHandler)
	SetUserChecker(namespace string, checker ChannelAuthChecker)
	HandleChannelAuth(w http.ResponseWriter, r *http.Request)
	HadleWebHook(w http.ResponseWriter, r *http.Request)
}

type ProviderConstructor func(ctx context.Context) ChannelProvider
