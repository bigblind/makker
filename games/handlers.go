package games

import (
	"net/http"
	"github.com/bigblind/makker/users"
	"github.com/bigblind/makker/di"
	"github.com/gorilla/mux"
	"github.com/bigblind/makker/handler_helpers"
)

func CreateInstace(w http.ResponseWriter, r *http.Request) {
	di.Graph.Invoke(func(constructor StoreConstructor) {
		store := constructor(r.Context())
		inter := GamesInteractor{store}

		uid := users.GetUserId(r)

		vars := mux.Vars(r)
		inst, err := inter.CreateInstance(vars["game"], uid)
		if err != nil {
			handler_helpers.RespondWithJSONError(w, 400, err)
		} else {
			handler_helpers.RespondWithJSON(w, 200, inst)
		}
	})
}