package games

import (
	"net/http"
	"github.com/bigblind/makker/users"
	"github.com/bigblind/makker/di"
	"github.com/gorilla/mux"
	"github.com/bigblind/makker/handler_helpers"
	"reflect"
)

func CreateInstace(w http.ResponseWriter, r *http.Request) {
	var storeConstructor StoreConstructor
	storeConstructor = di.Graph.ResolveByAssignableType(reflect.TypeOf(storeConstructor))[0].Interface().(StoreConstructor)
	store := storeConstructor(r.Context())
	inter := GamesInteractor{store}


	uid := users.GetUserId(r)

	vars := mux.Vars(r)
	inst, err := inter.CreateInstance(vars["game"], uid)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}