package recovermw_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang/tutorial/todo/internal/api/middleware/recovermw"
)

func TestRecover_PanickingHandlerReturns500JSON(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	recovermw.Recover(panicHandler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.True(t, strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json"))

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "500", body.Error.Code)
	assert.NotEmpty(t, body.Error.Message)
}

func TestRecover_NonPanickingHandlerDoesNotPanic(t *testing.T) {
	noPanicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	recovermw.Recover(noPanicHandler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
