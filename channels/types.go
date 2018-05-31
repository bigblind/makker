package channels

import (
	"net/http"
	"context"
)

type Channel interface{
	Public() bool
	Namespace() string
	Id() string
	ClientId() string
	Emit(event string, data interface{})
	EmitExcluding(socketId, event string, data interface{})
}

type ChannelProvider interface {
	NewChannel(ctx context.Context, namespace, id string, public bool) Channel
	OnJoin(namespace string, handler func(ctx context.Context, channel Channel, userId, socketId string))
	OnLeave(namespace string, handler func(ctx context.Context, channel Channel, userId, socketId string))
	SetUserChecker(namespace string, checker func(ctx context.Context, channel Channel, userId string) error)
	HandleChannelAuth(w http.ResponseWriter, r *http.Request)
	HadleWebHook(w http.ResponseWriter, r *http.Request)
}

type ProviderConstructor func(ctx context.Context) ChannelProvider

