package appengine

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"net/http"

	"github.com/bigblind/makker"
	_ "github.com/bigblind/makker/channels/pusher" // so Pusher gets injected
	"github.com/bigblind/makker/di"
	"context"
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
	di.Graph.Provide(NewGameStore)
	di.Graph.Provide(func() func(ctx context.Context) *http.Client {
		return func(ctx context.Context) *http.Client {
			return urlfetch.Client(ctx)
		}
	})
}
