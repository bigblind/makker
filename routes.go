package makker

import (
	"github.com/gorilla/mux"
	"github.com/bigblind/makker/users"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/users/me", users.MeHandler)

	r.Use(users.UserIDMiddleware)

	return r
}
