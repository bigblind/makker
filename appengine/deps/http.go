package deps

import (
	"github.com/bigblind/makker/di"
	"google.golang.org/appengine/urlfetch"
	"context"
	"net/http"
)

func init()  {
	di.Graph.Provide(func() func(ctx context.Context) *http.Client {
		return func(ctx context.Context) *http.Client {
			return urlfetch.Client(ctx)
		}
	})
}
