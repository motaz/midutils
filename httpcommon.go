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
