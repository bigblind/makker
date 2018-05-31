package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
	"github.com/bigblind/makker/games"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/channels"
	"github.com/bigblind/makker/channels/pusher"
)

func init() {
	router := makker.GetRouter()
	router.Use(contextMiddleware)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r.WithContext(appengine.NewContext(r)))
	}))

	appengine.Main()
}

func init() {
	di.Graph.Provide(func() games.StoreConstructor {
		return NewGameStore
	})

	di.Graph.Provide(func() channels.ProviderConstructor {
		return pusher.NewChannelProvider
	})

}
