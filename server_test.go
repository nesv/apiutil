package apiutil

import "log"

func ExampleNewHTTPSRedirectServer() {
	// Assume you have another *http.Server listening on port 443.
	srv := NewHTTPSRedirectServer(":http")
	if err := srv.ListenAndServe(); err != nil {
		log.Println("error starting http-to-https redirect server:", err)
	}
}
