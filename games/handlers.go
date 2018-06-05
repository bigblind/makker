package games

import (
	"encoding/json"
	"github.com/bigblind/makker/handler_helpers"
	"github.com/bigblind/makker/users"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func CreateInstace(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor(r.Context())

	uid := users.GetUserId(r)

	vars := mux.Vars(r)
	inst, err := inter.CreateInstance(vars["game"], uid)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}

func GetInstance(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor(r.Context())

	vars := mux.Vars(r)
	inst, err := inter.GetInstance(vars["instanceId"])
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor(r.Context())

	uid := users.GetUserId(r)
	vars := mux.Vars(r)
	err := inter.StartGame(vars["instanceId"], uid)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, map[string]bool{
			"success": true,
		})
	}

}

type MoveRequest struct {
	Move interface{}
}

func MakeMove(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor(r.Context())

	uid := users.GetUserId(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	var mr MoveRequest
	err = json.Unmarshal(body, &mr)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusBadRequest, err)
	}

	vars := mux.Vars(r)
	inst, err := inter.MakeMove(vars["instanceId"], uid, mr.Move)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}
