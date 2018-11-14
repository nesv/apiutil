package apiutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const jsonContentType = "application/json"

// WriteJSON marshals the given value v to its JSON representation, and writes
// it to an http.ResponseWriter with the given HTTP status code. This function
// also makes sure to set the "Content-Type" header to "application/json".
func WriteJSON(w http.ResponseWriter, v interface{}, status int) {
	p, err := json.Marshal(&v)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(status)
	w.Write(p)
}

// JSONError is similar to http.Error() in where it allows you to write an
// error to an http.ResponseWriter with a given HTTP status, except for that
// it will wrap your error in a JSON object, and put the error message under
// the object's "error" key.
//
// Calling this function will result in a response body like so:
//
// 	{"error": "...your error message..."}
//
func JSONError(w http.ResponseWriter, errStr string, status int) {
	WriteJSON(w, map[string]string{"error": errStr}, status)
}

// ReadJSON unmarshals the JSON data from the body of the given http.Request
// into v.
func ReadJSON(r *http.Request, v interface{}) error {
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(p, v)
}

// AcceptsJSON returns a boolean value indicating whether or not the Accept
// request header's value starts with "application/json".
func AcceptsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if len(accept) < len(jsonContentType) {
		return false
	}
	if accept[:len(jsonContentType)] == jsonContentType {
		return true
	}
	return false
}
