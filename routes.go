package makker

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/bigblind/makker/handler_helpers"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		handler_helpers.RespondWithJSON(w, 200, map[string]string{
			"hello": "world",
		})
	})


	return r
}
