package apiutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	input    interface{}
	want     []byte
	willFail bool
}

func TestWriteJSON(t *testing.T) {
	tests := []testCase{{
		input: map[string]uint64{"id": 123091241284012871},
		want:  []byte(`{"id":123091241284012871}`),
	}, {
		input: struct {
			Name string `json:"name"`
		}{
			Name: "nesv",
		},
		want: []byte(`{"name":"nesv"}`),
	},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		WriteJSON(rec, tc.input, http.StatusOK)

		if h := rec.Header().Get("Content-Type"); h != "application/json" {
			t.Errorf("Content-Type header has the incorrect value %q", h)
		}

		b := rec.Body.Bytes()
		if !bytes.Equal(tc.want, b) {
			t.Errorf("mismatched response bodies\nwanted %s\ngot %s", tc.want, b)
		}
	}
}

func TestJSONError(t *testing.T) {
	tests := map[string]string{
		"database error": `{"error":"database error"}`,
	}
	for in, want := range tests {
		rec := httptest.NewRecorder()
		JSONError(rec, in, http.StatusServiceUnavailable)

		if s := rec.Body.String(); s != want {
			t.Errorf("mismatched response bodies\nwanted %s\ngot %s", want, s)
		}
	}
}
