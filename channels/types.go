package channels

import "net/http"

type Channel interface{
	Namespace() string
	Id() string
	ClientId() string
	Emit(event string, data interface{})
	EmitExcluding(socketId, event string, data interface{})
}

type ChannelProvider interface {
	NewChannel(namespace, id string) Channel
	OnJoin(namespace string, handler func(channel Channel, userId, socketId string))
	OnLeave(namespace string, handler func(channel Channel, userId, socketId string))

	HandleChannelAuth(w http.ResponseWriter, r *http.Request)
	HadleWebHook(w http.ResponseWriter, r *http.Request)
}

