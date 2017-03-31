package apiutil

import "net/http"

// NewHTTPSRedirectServer returns an *http.Server that listens on addr for
// HTTP requests, and will return a response with the
// http.StatusMovedPermanently status code, redirecting clients to the same URL
// with the "https" scheme.
func NewHTTPSRedirectServer(addr string) *http.Server {
	return newServer(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		u.Scheme = "https"
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}))
}

func newServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: h,
	}
}
