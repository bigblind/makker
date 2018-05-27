package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
	"github.com/bigblind/makker/games"
	"github.com/bigblind/makker/di"
	"github.com/karlkfi/inject"
	"context"
)

func init() {
	r := makker.GetRouter()
	r.Use(contextMiddleware)
	http.Handle("/", r)

	appengine.Main()

	injectDeps()
}

func injectDeps() {
	var gs games.GameStore
	di.Graph.Define(&gs, inject.NewProvider(func() func(c context.Context) games.GameStore {
		return NewGameStore
	}))
}
