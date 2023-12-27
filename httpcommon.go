package midutils

import (
	"encoding/json"
	"net/http"
)

type ResponseType struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Errorcode int    `json:"errorcode"`
}

func SetStatusError(w http.ResponseWriter, message string, errorCode int,
	statusCode int) {

	w.WriteHeader(statusCode)

	w.Header().Set("content-type", "application/json")
	var res ResponseType
	res.Success = false
	res.Errorcode = errorCode
	res.Message = message
	data, _ := json.Marshal(res)
	w.Write(data)
	WriteLog("Error: " + res.Message)

}

func MethodGet(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CheckMethodAndContentType(w, r, "GET") {
			next.ServeHTTP(w, r)
		}
	}
}

func MethodPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CheckMethodAndContentType(w, r, "POST") {
			next.ServeHTTP(w, r)
		}
	}
}

func MethodDelete(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CheckMethodAndContentType(w, r, "DELETE") {
			next.ServeHTTP(w, r)
		}
	}
}
