package apiutil

import (
	"errors"
	"net/http"
	"strings"
)

var errAPIVersionRequired = errors.New("the X-API-Version header is required")

// VersionRouter is an http.Handler that allows you to route HTTP requests to
// a different http.Handler, based on the value of the X-API-Version request
// header.
//
// An empty string ("") as a key, is treated as the default handler to use,
// when no other routes match. If no default handler is specified,
// VersionRouter will respond with http.StatusBadRequest (400), with an
// error message indicating that the user must set the X-API-Version header in
// future requests.
//
// Making a request with the OPTIONS method will result in a response with the
// acceptable values for the X-API-Version header, separated by commas.
type VersionRouter map[string]http.Handler

// ServeHTTP implements the http.Handler interface.
func (h *VersionRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.Header().Set("X-API-Version-Required", "yes")

		versions := make([]string)
		for version := range h {
			if version == "" {
				continue
			}
			versions = append(versions, version)
		}
		w.Header().Set("X-API-Version", strings.Join(versions, ", "))

	default:
		handler, ok := h[r.Header.Get("X-API-Version")]
		if !ok {
			h.serveDefault(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func (h *VersionRouter) getDefaultHandler() http.Handler {
	if hh, ok := h[""]; ok {
		return hh
	}
	return nil
}

func (h *VersionRouter) serveDefault(w http.ResponseWriter, r *http.Request) {
	if hh := h.getDefaultHandler(); hh != nil {
		hh.ServeHTTP(w, r)
		return
	}

	// There is no default handler provided, so return an error message
	// indicating the X-API-Version header must be set on future requests.
	if AcceptsJSON(r) {
		JSONError(w, errAPIVersionRequired.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, errAPIVersionRequired.Error(), http.StatusBadRequest)
}
