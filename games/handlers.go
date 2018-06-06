package games

import (
	"encoding/json"
	"github.com/bigblind/makker/handler_helpers"
	"github.com/bigblind/makker/users"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func ListInstancesByGame(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor()

	vars := mux.Vars(r)

	insts, err := inter.ListInstances(r.Context(), vars["game"])
	if err != nil {
		handler_helpers.RespondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	handler_helpers.RespondWithJSON(w, http.StatusOK, insts)
}

func CreateInstace(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor()

	uid := users.GetUserId(r)

	vars := mux.Vars(r)
	inst, err := inter.CreateInstance(r.Context(), vars["game"], uid)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}

func GetInstance(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor()

	vars := mux.Vars(r)
	inst, err := inter.GetInstance(r.Context(), vars["instanceId"])
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	inter := NewInteractor()

	uid := users.GetUserId(r)
	vars := mux.Vars(r)
	err := inter.StartGame(r.Context(), vars["instanceId"], uid)
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
	inter := NewInteractor()

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
	inst, err := inter.MakeMove(r.Context(), vars["instanceId"], uid, mr.Move)
	if err != nil {
		handler_helpers.RespondWithJSONError(w, 400, err)
	} else {
		handler_helpers.RespondWithJSON(w, 200, inst)
	}
}
