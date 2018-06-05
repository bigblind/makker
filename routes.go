package makker

import (
	"github.com/bigblind/makker/channels"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/games"
	"github.com/bigblind/makker/users"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/users/me", users.MeHandler)

	r.HandleFunc("/games/{game}/new", games.CreateInstace).Methods("POST")
	r.HandleFunc("/games/{game}/instances", games.ListInstancesByGame).Methods("GET")
	r.HandleFunc("/games/instances/{instanceId}", games.GetInstance).Methods("GET")
	r.HandleFunc("/games/instances/{instanceId}/start", games.StartGame).Methods("POST")
	r.HandleFunc("/games/instances/{instanceId}/moves", games.MakeMove).Methods("POST")

	di.Graph.Invoke(func(cp channels.ChannelProvider) {
		r.HandleFunc("/channels/auth", cp.HandleChannelAuth).Methods("POST")
		r.HandleFunc("/channels/webhook", cp.HadleWebHook).Methods("POST")
	})

	r.Use(users.UserIdMiddleware)

	return r
}
