package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
	"github.com/bigblind/makker/games"
	"github.com/bigblind/makker/di"
	"github.com/karlkfi/inject"
)

func init() {
	router := makker.GetRouter()
	router.Use(contextMiddleware)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r.WithContext(appengine.NewContext(r)))
	}))

	injectDeps()
	appengine.Main()
}

func injectDeps() {
	var gs games.StoreConstructor
	di.Graph.Define(&gs, inject.NewProvider(func() games.StoreConstructor {
		return NewGameStore
	}))
}
