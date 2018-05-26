package appengine

import (
	"net/http"
	"google.golang.org/appengine"
)

type contextHandler struct {
	wrapped http.Handler
}

func (ch contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	r = r.WithContext(ctx)
	ch.wrapped.ServeHTTP(w, r)
}

func contextMiddleware(h http.Handler) http.Handler {
	return contextHandler{h}
}