package httphelper

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"app_api/shared"
	loghelper "app_api/shared/loghelpers"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Error   int         `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type BatchResponse struct {
	Result interface{} `json:"result"`
	Length int64       `json:"count"`
	Offset int64       `json:"offset"`
	Limit  int64       `json:"limit"`
	Total  int64       `json:"total"`
}

// NewResponse give a response
// If error is empty status code will be 200 OK; otherwise code will be retrieved from apiError
// If error consist code 200 - method will return 200 OK, not an error
func NewResponse(ctx context.Context, w http.ResponseWriter, result interface{}, apiError *shared.APIError) {
	var resp response
	w.Header().Set("content-type", "application/json")
	if apiError == nil || (apiError.ErrorCode >= 200 && apiError.ErrorCode <= 299) {
		// Currently successful calls all respond with the same Status Code and message. We should ticket this for future refactor
		log.WithContext(ctx).WithFields(log.Fields{
			"Response":     result,
			"ResponseCode": http.StatusOK,
			"RequestID":    loghelper.GetRequestID(ctx),
			"Request":      loghelper.GetRequestPath(ctx),
		}).Info()

		jResult, err := json.Marshal(result)
		if err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"ErrorMessage": err,
				"RequestID":    loghelper.GetRequestID(ctx),
				"Request":      loghelper.GetRequestPath(ctx),
			}).Error()
			w.WriteHeader(resp.statusCode(http.StatusInternalServerError))
			jResult = []byte(`{"error": 1, "message": "error"}`)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jResult); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"ErrorMessage": err,
				"RequestID":    loghelper.GetRequestID(ctx),
				"Request":      loghelper.GetRequestPath(ctx),
			}).Error()
		}
	} else {
		w.WriteHeader(resp.statusCode(apiError.ErrorCode))
		loghelper.LogErrors(ctx, apiError)
		resp = response{
			Error:   1,
			Message: apiError.ClientErrorMessage,
		}
		if apiError.ValidationErrorMessage != "" {
			resp.Result = map[string]interface{}{"validation_error": strings.Split(apiError.ValidationErrorMessage, "\n")}
		}

		jResult, err := resp.marshal()
		if err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"ErrorMessage": err,
				"RequestID":    loghelper.GetRequestID(ctx),
				"Request":      loghelper.GetRequestPath(ctx),
			}).Error()
			w.WriteHeader(resp.statusCode(http.StatusInternalServerError))
			jResult = []byte(`{"error": 1, "message": "error"}`)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jResult); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"ErrorMessage": err,
				"RequestID":    loghelper.GetRequestID(ctx),
				"Request":      loghelper.GetRequestPath(ctx),
			}).Error()
		}
	}
}

// prevent a panic if status code will be wrong
func (r response) statusCode(code int) int {
	if code < 100 || code > 999 {
		return http.StatusInternalServerError
	}
	return code
}

func (r response) marshal() ([]byte, error) {
	return json.Marshal(r)
}
