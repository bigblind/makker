package appengine

import (
	"google.golang.org/appengine/log"
	"net/http"
	"runtime/debug"
)

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		debug.SetTraceback("all")
		defer func() {
			if rv := recover(); rv != nil {
				log.Errorf(ctx, "Application panicnked!\n%v\n%s", rv, debug.Stack())
			}
		}()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
