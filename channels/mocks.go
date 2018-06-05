package channels

import (
	"github.com/stretchr/testify/mock"
	"context"
	"net/http"
)

type MockChannelProvider struct {
	mock.Mock
}

func (mcp *MockChannelProvider) NewChannel(ctx context.Context, namespace, id string, public bool) Channel {
	vals := mcp.Called(ctx, namespace, id, public)
	return vals.Get(0).(Channel)
}

func (mcp *MockChannelProvider) OnJoin(namespace string, handler EventHandler) {
	mcp.Called(namespace, handler)
}

func (mcp *MockChannelProvider) OnLeave(namespace string, handler EventHandler) {
	mcp.Called(namespace, handler)
}

func (mcp *MockChannelProvider) SetUserChecker(namespace string, checker ChannelAuthChecker) {
	mcp.Called(namespace, checker)
}

func (mcp *MockChannelProvider) HandleChannelAuth(w http.ResponseWriter, r *http.Request) {
	mcp.Called(w, r)
}

func (mcp *MockChannelProvider) HadleWebHook(w http.ResponseWriter, r *http.Request) {
	mcp.Called(w, r)
}

type MockChannel struct {
	mock.Mock
	namespace, id string
	public bool
}

func NewMockChannel(namespace, id string, public bool) *MockChannel {
	mc := new(MockChannel)
	mc.namespace = namespace
	mc.id = id
	mc.public = public
	return mc
}

func (mc *MockChannel) Public() bool {
	return mc.public
}

func (mc *MockChannel) Namespace() string {
	return mc.namespace
}

func (mc *MockChannel) Id() string {
	return mc.id
}

func (mc *MockChannel) ClientId() string {
	args := mc.Called()
	return args.String(0)
}

func (mc *MockChannel) Emit(event string, data interface{}) {
	mc.Called(event, data)
}

func (mc *MockChannel) EmitExcluding(socketId, event string, data interface{}) {
	mc.Called(socketId, event, data)
}


