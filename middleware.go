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
		/** AssignRequestID ... is just a UUID to trace requests back, which helps with debugging.
		But by attaching it to the individual requests' context allows it to become traceable throughout the API.

		e.g. The following are the initial request log, and future error log:
		time="2020-11-08T16:25:06-08:00" level=info Method=POST RemoteAddress="[::1]:59276" RequestID=9eb198d6-91b7-46f6-b8f5-83ee71351b3f URL=/vehicles/1234/engine
		time="2020-11-08T16:25:06-08:00" level=error Caller="Called from app_api/apis/vehicle.(*service).SendEngineAction; /Users/andyqi/Desktop/Projects/go/src/SmartCar/apis/vehicle/vehicle.go#147" ClientError="Unsupported engine action option" Error="Unsupported data type: FOOBAR" ErrorCode=500 InternalError="Unsupported engine action option" Request=POST./vehicles/1234/engine RequestID=9eb198d6-91b7-46f6-b8f5-83ee71351b3f
		*/
		ctx := loghelper.AssignRequestID(r.Context())
		ctx = loghelper.AssignRequestPath(ctx, r)
		r = r.WithContext(ctx)

		// Example log:
		// time="2020-11-08T16:21:12-08:00" level=info Method=GET RemoteAddress="[::1]:59201" RequestID=48497053-376c-4bd3-916b-26d2595e2379 URL=/vehicles/1234/battery
		log.WithContext(ctx).WithFields(log.Fields{
			"Method":        r.Method,
			"URL":           r.RequestURI,
			"RemoteAddress": r.RemoteAddr,
			"RequestID":     loghelper.GetRequestID(r.Context()),
		}).Info()

		next.ServeHTTP(w, r)
	})
}
