package main

import (
	"net/http"

	loghelper "app_api/shared/loghelpers"

	log "github.com/sirupsen/logrus"
)

// ContextKey ... type definition
type ContextKey string

// Logger ... logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := loghelper.AssignRequestID(r.Context())
		ctx = loghelper.AssignRequestPath(ctx, r)
		r = r.WithContext(ctx)
		log.WithContext(ctx).WithFields(log.Fields{
			"Method":        r.Method,
			"URL":           r.RequestURI,
			"RemoteAddress": r.RemoteAddr,
			"RequestID":     loghelper.GetRequestID(r.Context()),
		}).Info()

		next.ServeHTTP(w, r)
	})
}
