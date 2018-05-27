package makker

import (
	"github.com/gorilla/mux"
	"github.com/bigblind/makker/users"
	"github.com/bigblind/makker/games"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/users/me", users.MeHandler)

	r.HandleFunc("/games/{game}/new", games.CreateInstace).Methods("POST")

	r.Use(users.UserIdMiddleware)

	return r
}
