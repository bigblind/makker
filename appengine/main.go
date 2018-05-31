package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
	"github.com/bigblind/makker/games"
	"github.com/bigblind/makker/di"
	_ "github.com/bigblind/makker/channels/pusher" // so Pusher gets injected
)

func init() {
	initDeps()
	router := makker.GetRouter()
	router.Use(contextMiddleware)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r.WithContext(appengine.NewContext(r)))
	}))

	appengine.Main()
}

func initDeps() {
	di.Graph.Provide(func() games.StoreConstructor {
		return NewGameStore
	})
}
