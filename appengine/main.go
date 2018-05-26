package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
)

func init() {
	r := makker.GetRouter()
	r.Use(contextMiddleware)
	http.Handle("/", r)
	appengine.Main()
}
