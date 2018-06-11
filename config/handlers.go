package config

import (
	"net/http"
	"github.com/bigblind/makker/handler_helpers"
)

func GetConfig(w http.ResponseWriter, r *http.Request)  {
	handler_helpers.RespondWithJSON(w, 200, map[string]string{
		"pusher_key": PusherKey,
		"pusher_cluster": PusherCluster,
	})
}
