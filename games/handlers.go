package games

import (
	"net/http"
	"github.com/bigblind/makker/users"
	"github.com/bigblind/makker/di"
	"github.com/gorilla/mux"
	"github.com/bigblind/makker/handler_helpers"
	"github.com/bigblind/makker/channels"
)

func CreateInstace(w http.ResponseWriter, r *http.Request) {
	di.Graph.Invoke(func(constructor StoreConstructor, providerConstructor channels.ProviderConstructor) {
		inter := NewInteractor(r.Context())

		uid := users.GetUserId(r)

		vars := mux.Vars(r)
		inst, err := inter.CreateInstance(vars["game"], uid)
		if err != nil {
			handler_helpers.RespondWithJSONError(w, 400, err)
		} else {
			handler_helpers.RespondWithJSON(w, 200, instancetoResponse(&inst, r, providerConstructor(r.Context())))
		}
	})
}

type instanceResponse struct{
	*GameInstance
	PublicChannel  string
	PrivateChannel string
}

func instancetoResponse(i *GameInstance, r *http.Request, cp channels.ChannelProvider) instanceResponse {
	uid := users.GetUserId(r)
	chanIds := i.Channels(uid)
	return instanceResponse{
		GameInstance: i,
		PublicChannel: cp.NewChannel(r.Context(), "games", chanIds.Public, true).ClientId(),
		PrivateChannel: cp.NewChannel(r.Context(), "games", chanIds.Private, false).ClientId(),
	}
}