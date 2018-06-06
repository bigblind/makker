package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	// set up deps
	_ "deps"
	_ "github.com/bigblind/makker/channels/pusher"

	"github.com/bigblind/makker"
)

func init() {
	router := makker.GetRouter()
	router.Use(contextMiddleware)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r.WithContext(appengine.NewContext(r)))
	}))

	appengine.Main()
}
