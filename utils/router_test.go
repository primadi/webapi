package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	router := NewRouter("")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestInfo(t *testing.T) {
	router := NewRouter("")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	rm := jsoniter.Get(w.Body.Bytes(), "runMode").ToString()
	assert.True(t, rm == "debug" || rm == "release")
}
