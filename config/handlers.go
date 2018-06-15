package config

import (
	"net/http"
	"github.com/bigblind/makker/handler_helpers"
)

var FrontendConfig = map[string]string{
	"pusher_key": PusherKey,
	"pusher_cluster": PusherCluster,
}

func GetConfig(w http.ResponseWriter, r *http.Request)  {
	handler_helpers.RespondWithJSON(w, 200, FrontendConfig)
}
