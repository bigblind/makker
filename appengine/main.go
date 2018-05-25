package appengine

import (
	"google.golang.org/appengine"
	"net/http"

	"github.com/bigblind/makker"
)

func init() {
	http.Handle("/", makker.GetRouter())
	appengine.Main()
}
