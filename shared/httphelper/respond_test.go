package httphelper

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"app_api/shared"

	"github.com/stretchr/testify/assert"
)

const (
	contentType = "content-type"
)

func TestNewResponseSuccess(t *testing.T) {
	expectedResponse := response{
		Message: "success",
	}
	jExpected, err := json.Marshal(expectedResponse)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	ctx := context.Background()
	NewResponse(ctx, w, nil, nil)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "application/json", resp.Header.Get(contentType))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, jExpected, body)
}

func TestNewResponseSuccessWithResult(t *testing.T) {
	w := httptest.NewRecorder()
	type RespStruct struct {
		Username string
		Address  string
		Code     int
	}
	res := RespStruct{"Test", "Address", 123}
	ctx := context.Background()
	NewResponse(ctx, w, res, nil)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "application/json", resp.Header.Get(contentType))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"error":0,"message":"success","result":{"Username":"Test","Address":"Address","Code":123}}`, string(body))
}

func TestNewResponseError(t *testing.T) {
	expectedResponse := response{
		Error:   1,
		Message: "Internal Error",
	}
	jExpected, err := json.Marshal(expectedResponse)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	msg := "error"
	e := shared.NewAPIError(http.StatusBadRequest, errors.New(msg), "Internal Error").
		SetInternalErrorMessage(msg)
	ctx := context.Background()
	NewResponse(ctx, w, nil, e)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "application/json", resp.Header.Get(contentType))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, jExpected, body)
}
