package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
	_ "github.com/bigblind/makker/channels/pusher" // so Pusher gets injected
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/games"
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
