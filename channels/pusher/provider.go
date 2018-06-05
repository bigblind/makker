package pusher

import (
	"context"
	"github.com/bigblind/makker/channels"
	"github.com/bigblind/makker/config"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/handler_helpers"
	"github.com/bigblind/makker/users"
	"github.com/pusher/pusher-http-go"
	"go.uber.org/dig"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	di.Graph.Provide(newChannelProvider)
}

type PusherProvider struct {
	HttpClientConstructor func(ctx context.Context) *http.Client

	JoinListeners  map[string][]channels.EventHandler
	LeaveListeners map[string][]channels.EventHandler
	UserCheckers   map[string]channels.ChannelAuthChecker
}

type channelProviderParams struct {
	dig.In

	ClientConstructor func(ctx context.Context) *http.Client `optional:"true"`
}

func newChannelProvider(params channelProviderParams) channels.ChannelProvider {
	var pp PusherProvider
	pp = PusherProvider{
		HttpClientConstructor: params.ClientConstructor,
		JoinListeners:         make(map[string][]channels.EventHandler),
		LeaveListeners:        make(map[string][]channels.EventHandler),
		UserCheckers:          make(map[string]channels.ChannelAuthChecker),
	}

	return pp
}

func (pp PusherProvider) client(ctx context.Context) pusher.Client {
	c := pusher.Client{
		AppId:   config.PusherAppId,
		Key:     config.PusherKey,
		Secret:  config.PusherSecret,
		Cluster: config.PusherCluster,
		Secure:  true,
	}

	if pp.HttpClientConstructor != nil {
		c.HttpClient = pp.HttpClientConstructor(ctx)
	}

	return c
}

func (pp PusherProvider) NewChannel(ctx context.Context, namespace, id string, public bool) channels.Channel {
	c := pp.client(ctx)
	pc := PusherChannel{
		client:    &c,
		namespace: namespace,
		id:        id,
		public:    public,
	}

	return &pc
}

func (pp PusherProvider) ChannelFromClientId(ctx context.Context, id string) channels.Channel {
	// Public channels have the format
	//    {namespace}-{id}
	//
	// Private channels have the format
	//    presence-{namespace}-{id}
	//

	parts := strings.SplitN(id, "-", 3)
	if parts[0] == "presence" {
		return pp.NewChannel(ctx, parts[1], parts[2], false)
	} else {
		return pp.NewChannel(ctx, parts[0], parts[1], true)
	}
}

func (pp PusherProvider) OnJoin(namespace string, handler channels.EventHandler) {
	var listeners []channels.EventHandler
	var ok bool

	if listeners, ok = pp.JoinListeners[namespace]; !ok {
		listeners = make([]channels.EventHandler, 1)
	}

	listeners = append(listeners, handler)
	pp.JoinListeners[namespace] = listeners
}

func (pp PusherProvider) OnLeave(namespace string, handler channels.EventHandler) {
	var listeners []channels.EventHandler
	var ok bool

	if listeners, ok = pp.LeaveListeners[namespace]; !ok {
		listeners = make([]channels.EventHandler, 1)
	}

	listeners = append(listeners, handler)
	pp.LeaveListeners[namespace] = listeners
}

func (pp PusherProvider) SetUserChecker(namespace string, checker channels.ChannelAuthChecker) {
	pp.UserCheckers[namespace] = checker
}

func (pp PusherProvider) HandleChannelAuth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	uid := users.GetUserId(r)
	m := pusher.MemberData{
		UserId:   uid,
		UserInfo: make(map[string]string),
	}

	c := pp.client(r.Context())
	resp, err := c.AuthenticatePresenceChannel(body, m)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusForbidden, err)
		return
	}

	// We can safely ignore the error here, because AuthenticatePresenceChannel does the same check.
	params, _ := url.ParseQuery(string(body))
	ch := pp.ChannelFromClientId(r.Context(), params["channel_name"][0])

	if checker, ok := pp.UserCheckers[ch.Namespace()]; ok {
		err = checker(r.Context(), ch, uid)
		if err != nil {
			handler_helpers.RespondWithJSONError(w, http.StatusForbidden, err)
		}
	}

	w.Write(resp)
}

func (pp PusherProvider) HadleWebHook(w http.ResponseWriter, r *http.Request) {
	c := pp.client(r.Context())
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	wh, err := c.Webhook(r.Header, body)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	var listeners []channels.EventHandler
	var ok bool

	for _, e := range wh.Events {
		ch := pp.ChannelFromClientId(r.Context(), e.Channel)

		switch e.Name {
		case "member_added":
			listeners, ok = pp.JoinListeners[ch.Namespace()]
			break
		case "member_removed":
			listeners, ok = pp.LeaveListeners[ch.Namespace()]
			break
		default:
			continue
		}

		if !ok {
			continue
		}

		for _, l := range listeners {
			l(r.Context(), ch, e.UserId, e.SocketId)
		}
	}
}
