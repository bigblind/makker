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



var store = sessions.NewCookieStore(config.Secret)

func UserIdMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")

		_, ok := session.Values["userId"]
		if !ok {
			session.Values["userId"] = fmt.Sprintf("guest_player%v", rand.Int())
		}
		session.Save(r, w)

		next.ServeHTTP(w, r)
	})
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