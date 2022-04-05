package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() (fiboHandler, *httptest.ResponseRecorder, *http.Request)
		assert func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			it: `given an empty fibo handler
				when serving the "/" endpoint with get method
				then return message "Server is Happy!"`,
			setup: func() (fiboHandler, *httptest.ResponseRecorder, *http.Request) {
				fh := fiboHandler{}
				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				return fh, w, req
			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				res := w.Body.String()
				assert.Equal(t, "Server is happy!", res)
			},
		},
		{
			it: `given a fibo handler with prev 3 and current 5
				when serving the "/previous" endpoint with get method
				then return message "3"`,
			setup: func() (fiboHandler, *httptest.ResponseRecorder, *http.Request) {
				fh := fiboHandler{
					prev:    3,
					current: 5,
				}
				req := httptest.NewRequest("GET", "/previous", nil)
				w := httptest.NewRecorder()
				return fh, w, req
			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				res := w.Body.String()
				assert.Equal(t, "3", res)
			},
		},
		{
			it: `given a fibo handler with prev 3 and current 5
				when serving the "/current" endpoint with get method
				then return message "5"`,
			setup: func() (fiboHandler, *httptest.ResponseRecorder, *http.Request) {
				fh := fiboHandler{
					prev:    3,
					current: 5,
				}
				req := httptest.NewRequest("GET", "/current", nil)
				w := httptest.NewRecorder()
				return fh, w, req
			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				res := w.Body.String()
				assert.Equal(t, "5", res)
			},
		},
		{
			it: `given a fibo handler with prev 3 and current 5
				when serving the "/next" endpoint with get method
				then return message "8"`,
			setup: func() (fiboHandler, *httptest.ResponseRecorder, *http.Request) {
				fh := fiboHandler{
					prev:    3,
					current: 5,
				}
				req := httptest.NewRequest("GET", "/next", nil)
				w := httptest.NewRecorder()
				return fh, w, req
			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				res := w.Body.String()
				assert.Equal(t, "8", res)
			},
		},
		{
			it: `given an empty fibo handler
				when serving the a random endpoint with get method
				then return message "not-found"`,
			setup: func() (fiboHandler, *httptest.ResponseRecorder, *http.Request) {
				fh := fiboHandler{}
				req := httptest.NewRequest("GET", "/random", nil)
				w := httptest.NewRecorder()
				return fh, w, req
			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				res := w.Body.String()
				assert.Equal(t, "not-found", res)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			fh, w, req := tt.setup()
			fh.ServeHTTP(w, req)
			tt.assert(t, w)
		})
	}
}

func TestServeHTTPWithPanic(t *testing.T) {
	fh := fiboHandler{}
	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	assert.Panics(t, func() {
		fh.ServeHTTP(w, req)
	})
}

func TestHandlerRoot(t *testing.T) {
	assert.HTTPBodyContains(t, handlerRoot, "GET", "/", nil, "Server is happy!")
}

func TestHandlerPrevious(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() fiboHandler
		assert func(*testing.T, fiboHandler)
	}{
		{
			it: `given prev equals 3
				and current equals 5
				when hitting the /previous endpoint
				then returns 3`,
			setup: func() fiboHandler {
				return fiboHandler{
					prev:    3,
					current: 5,
				}
			},
			assert: func(t *testing.T, fh fiboHandler) {
				assert.HTTPBodyContains(t, fh.handlerPrevious, "GET", "/previous", nil, "3")
			},
		},
		{
			it: `given prev equals 21
				and current equals 34
				when hitting the /previous endpoint
				then returns 21`,
			setup: func() fiboHandler {
				return fiboHandler{
					prev:    21,
					current: 34,
				}
			},
			assert: func(t *testing.T, fh fiboHandler) {
				assert.HTTPBodyContains(t, fh.handlerPrevious, "GET", "/previous", nil, "21")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			handler := tt.setup()
			tt.assert(t, handler)
		})
	}
}

func TestHandlerCurrent(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() fiboHandler
		assert func(*testing.T, fiboHandler)
	}{
		{
			it: `given prev equals 3
				and current equals 5
				when hitting the /current endpoint
				then returns 5`,
			setup: func() fiboHandler {
				return fiboHandler{
					prev:    3,
					current: 5,
				}
			},
			assert: func(t *testing.T, fh fiboHandler) {
				assert.HTTPBodyContains(t, fh.handlerCurrent, "GET", "/current", nil, "5")
			},
		},
		{
			it: `given prev equals 21
				and current equals 34
				when hitting the /previous endpoint
				then returns 34`,
			setup: func() fiboHandler {
				return fiboHandler{
					prev:    21,
					current: 34,
				}
			},
			assert: func(t *testing.T, fh fiboHandler) {
				assert.HTTPBodyContains(t, fh.handlerCurrent, "GET", "/current", nil, "34")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			handler := tt.setup()
			tt.assert(t, handler)
		})
	}
}

func TestHandlerNext(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() fiboHandler
		assert func(*testing.T, fiboHandler)
	}{
		{
			it: `given prev equals 3
			and current equals 5
			when hitting the /current endpoint
			then returns 8`,
			setup: func() fiboHandler {
				return fiboHandler{
					prev:    3,
					current: 5,
				}
			},
			assert: func(t *testing.T, fh fiboHandler) {
				assert.HTTPBodyContains(t, fh.handlerNext, "GET", "/next", nil, "8")
			},
		},
		{
			it: `given prev equals 21
			and current equals 34
			when hitting the /previous endpoint
			then returns 55`,
			setup: func() fiboHandler {
				return fiboHandler{
					prev:    21,
					current: 34,
				}
			},
			assert: func(t *testing.T, fh fiboHandler) {
				assert.HTTPBodyContains(t, fh.handlerNext, "GET", "/next", nil, "55")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			handler := tt.setup()
			tt.assert(t, handler)
		})
	}
}

func TestHandlerNotFound(t *testing.T) {
	assert.HTTPBodyContains(t, handlerNotFound, "GET", "/random", nil, "not-found")
}

func TestHandlerError(t *testing.T) {
	assert.Panics(t, func() {
		handlerError(nil, nil)
	})
}
