package users

import (
	"github.com/gorilla/sessions"
	"github.com/bigblind/makker/config"
	"net/http"
	"fmt"
	"math/rand"
	"strings"
	"github.com/bigblind/makker/handler_helpers"
)

type userIdHandler struct {
	wrappedHandler http.Handler
}

var store = sessions.NewCookieStore(config.Secret)

func (h userIdHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")

	_, ok := session.Values["userId"]
	if !ok {
		session.Values["userId"] = fmt.Sprintf("guest_player%v", rand.Int())
	}

	h.wrappedHandler.ServeHTTP(w, r)

	session.Save(r, w)
}

func UserIDMiddleware(handler http.Handler) http.Handler {
	return userIdHandler{handler}
}

func GetUserId(r *http.Request) string {
	session, _ := store.Get(r, "auth")

	id, ok := session.Values["userId"]
	if !ok {
		return ""
	}

	return id.(string)
}


type UserData struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func MeHandler (w http.ResponseWriter, r *http.Request) {
	id := GetUserId(r)
	name := strings.Split(id, "_")[1]
	handler_helpers.RespondWithJSON(w, 200, UserData{id, name})
}