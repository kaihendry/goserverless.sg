package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecorder(t *testing.T) {
	type checkFunc func(*httptest.ResponseRecorder) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasStatus := func(want int) checkFunc {
		return func(rec *httptest.ResponseRecorder) error {
			if rec.Code != want {
				return fmt.Errorf("expected status %d, found %d", want, rec.Code)
			}
			return nil
		}
	}
	containsContents := func(want string) checkFunc {
		return func(rec *httptest.ResponseRecorder) error {
			if have := rec.Body.String(); !strings.Contains(have, want) {
				return fmt.Errorf("expected to find %q, in %q", want, have)
			}
			return nil
		}
	}
	hasHeader := func(key, want string) checkFunc {
		return func(rec *httptest.ResponseRecorder) error {
			if have := rec.Result().Header.Get(key); have != want {
				return fmt.Errorf("expected header %s: %q, found %q", key, want, have)
			}
			return nil
		}
	}

	tests := [...]struct {
		name   string
		h      func(w http.ResponseWriter, r *http.Request)
		checks []checkFunc
	}{
		{
			"200 default",
			handleIndex,
			check(hasStatus(200), containsContents("DOCTYPE"), hasHeader("Content-Type", "text/html; charset=utf-8")),
		},
		{
			"200 rank",
			handleRank,
			check(hasStatus(200), containsContents("ap-southeast-1"), hasHeader("Content-Type", "application/json")),
		},
	}

	r, _ := http.NewRequest("GET", "https://goserverless.sg/", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(tt.h)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, r)
			for _, check := range tt.checks {
				if err := check(rec); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
