package httphelper

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"app_api/shared"
	testhelper "app_api/shared/testhelpers"

	"github.com/stretchr/testify/assert"
)

type testJSON struct {
	Integer int     `json:"integer"`
	String  string  `json:"string"`
	Float   float64 `json:"float"`
	Bool    bool    `json:"bool"`
}

const good = `{"integer": 1, "string": "characters", "float": 1.5, "bool": true}`
const badForm = `{"integer":1,"string":"characters","float":1.5,"bool" true}`
const unexpEOF = `{`
const invalidVal = `{"integer":1,"string":123,"float":1.5,"bool":"hello"}`
const unknownField = `{"integer":1,"unknown":"field","string":"characters","float":1.5,"bool":true}`
const empty = ``
const multiple = `{"integer":1,"string":"characters","float":1.5,"bool":true}{"integer":1,"string":"characters","float":1.5,"bool":true}`
const goodCType = "application/json"
const badCType = "application/x-www-form-urlencoded"

// HELPER FUNCTIONS
func jsonTestHelper(jsonString string, cType string) *shared.APIError {
	body := testJSON{}

	jsb := bytes.NewBuffer([]byte(jsonString))

	r, _ := http.NewRequest("POST", "/v3/account", jsb)
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", cType)

	err := DecodeJSONBody(w, r, &body)
	return err
}

func makeJSONTooLarge() string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	stringWithCharset := func(length int, charset string) string {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
		return string(b)
	}

	init := `{"integer":1,"string":"characters","float":1.5,"bool":true`
	round := 1

	for round < 40500 {
		init = init + `,"` + stringWithCharset(10, charset) + `":"` +
			stringWithCharset(10, charset) + `"`
		if round == 40499 {
			init = init + `}`
		}
		round++
	}
	return init
}

// TESTS
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDecodeJSONBodySuccess(t *testing.T) {
	err := jsonTestHelper(good, goodCType)

	assert.Nil(t, err)
}

// the following are all DecodeJSONBody fail modes
func TestContentType(t *testing.T) {
	err := jsonTestHelper(good, badCType)

	testhelper.CheckResponseCode(t, http.StatusUnsupportedMediaType, err.ErrorCode)
}

func TestBadForm(t *testing.T) {
	err := jsonTestHelper(badForm, goodCType)

	testhelper.CheckResponseCode(t, http.StatusBadRequest, err.ErrorCode)
}

func TestUnexpEOF(t *testing.T) {
	err := jsonTestHelper(unexpEOF, goodCType)

	testhelper.CheckResponseCode(t, http.StatusBadRequest, err.ErrorCode)
}

func TestInvalidVal(t *testing.T) {
	err := jsonTestHelper(invalidVal, goodCType)

	testhelper.CheckResponseCode(t, http.StatusBadRequest, err.ErrorCode)
}

func TestUnknownField(t *testing.T) {
	err := jsonTestHelper(unknownField, goodCType)

	testhelper.CheckResponseCode(t, http.StatusBadRequest, err.ErrorCode)
}

func TestEmpty(t *testing.T) {
	err := jsonTestHelper(empty, goodCType)

	testhelper.CheckResponseCode(t, http.StatusBadRequest, err.ErrorCode)
}

func TestTooLarge(t *testing.T) {
	tooLarge := makeJSONTooLarge()
	err := jsonTestHelper(tooLarge, goodCType)

	testhelper.CheckResponseCode(t, http.StatusRequestEntityTooLarge, err.ErrorCode)
}

func TestMultiple(t *testing.T) {
	err := jsonTestHelper(multiple, goodCType)

	testhelper.CheckResponseCode(t, http.StatusBadRequest, err.ErrorCode)
}
