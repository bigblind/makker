package pusher

import (
	"fmt"
	"github.com/pusher/pusher-http-go"
)

type PusherChannel struct {
	namespace, id string
	public        bool

	client *pusher.Client
}

func (pc *PusherChannel) Public() bool {
	return pc.public
}

func (pc *PusherChannel) Namespace() string {
	return pc.namespace
}

func (pc *PusherChannel) Id() string {
	return pc.id
}

func (pc *PusherChannel) ClientId() string {
	prefix := ""
	if !pc.public {
		prefix = "presence-"
	}

	return fmt.Sprintf("%v%v-%v", prefix, pc.namespace, pc.id)
}

func (pc *PusherChannel) Emit(event string, data interface{}) {
	pc.client.Trigger(pc.ClientId(), event, data)
}

func (pc *PusherChannel) EmitExcluding(socketId, event string, data interface{}) {
	pc.client.TriggerExclusive(pc.ClientId(), event, data, socketId)
}
