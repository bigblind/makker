package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
	"github.com/bigblind/makker/games"
	"github.com/bigblind/makker/di"
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
	di.Graph.Provide(func() games.StoreConstructor {
		return NewGameStore
	})
}
