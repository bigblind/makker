package deps

import (
	"context"
	"github.com/bigblind/makker/di"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func init() {
	di.Graph.Provide(func() func(context.Context) *http.Client {
		return func(ctx context.Context) *http.Client {
			return urlfetch.Client(ctx)
		}
	})
}
