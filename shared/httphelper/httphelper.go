package httphelper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"app_api/shared"

	"github.com/golang/gddo/httputil/header"
)

// HandleJSONErrors ... putting these checks in a separate function
// so that they can be used for any JSON checks, not only HTTP bodies
func HandleJSONErrors(ctx context.Context, body io.Reader, dst interface{}) (*json.Decoder, *shared.APIError) {
	// decode the json body and don't allow unexpected fields
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), "Request body contains badly-formed JSON").
				SetInternalErrorMessage(msg)
			return nil, e

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), "Request body contains badly-formed JSON").
				SetInternalErrorMessage(msg)
			return nil, e

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), msg)
			return nil, e

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), msg)
			return nil, e

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), msg)
			return nil, e

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			e := shared.NewAPIError(http.StatusRequestEntityTooLarge, errors.New(msg), msg)
			return nil, e

		default:
			e := shared.NewAPIError(http.StatusInternalServerError, err, "Request body contains badly-formed JSON").
				SetInternalErrorMessage("Failed to decode JSON")
			return nil, e
		}
	}

	return dec, nil
}

// DecodeJSONBody ... decode json body and perform various checks
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *shared.APIError {
	ctx := r.Context()

	// check for expected content-type header
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			e := shared.NewAPIError(http.StatusUnsupportedMediaType, errors.New(msg), "Content-Type header is not application/json").
				SetInternalErrorMessage(msg)
			return e
		}
	}

	// limit body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec, e := HandleJSONErrors(ctx, r.Body, dst)
	if e != nil {
		return e
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), msg)
		return e
	}
	return nil
}
