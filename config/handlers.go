package config

import (
	"net/http"
	"github.com/bigblind/makker/handler_helpers"
	"github.com/bigblind/makker/games/interactors"
)

func GetConfig(w http.ResponseWriter, r *http.Request)  {
	handler_helpers.RespondWithJSON(w, 200, map[string]string{
		"pusher_key": PusherKey,
		"pusher_cluster": PusherCluster,

		"lobby_channel": interactors.Interactor.LobbyChannel(r.Context()).ClientId(),
	})
}
