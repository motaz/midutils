package midutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/motaz/codeutils"
)

var logTableName string

func SetLogTableName(tablename string) {
	logTableName = tablename
}

func callLogger(req LogRequest, r *http.Request) {

	token := r.Header.Get("token")
	_, session, _ := CheckToken(token)
	req.Username = session.Session.Username
	req.IP = codeutils.GetRemoteIP(r)

	aurl := GetConfigValue("logger_api", "http://localhost:7654/logger/insert")
	reqData, _ := json.Marshal(req)
	var reqMap = make(map[string]interface{})
	json.Unmarshal(reqData, &reqMap)
	callURLPost(aurl, reqMap)
}

type LogRequest struct {
	TableName    string      `json:"table_name"`
	IP           string      `json:"ip"`
	Username     string      `json:"username"`
	MDN          string      `json:"mdn"`
	Errorcode    int         `json:"errorcode"`
	Message      string      `json:"message"`
	IsSuccessful bool        `json:"is_successful"`
	Request      interface{} `json:"request"`
	Response     interface{} `json:"response"`
}

type LogResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Log(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var tableName = logTableName + strings.ReplaceAll(r.URL.Path, "/", "_")

		var mdn string
		log := LogRequest{TableName: tableName}
		if r.Method == http.MethodGet {
			mdn = r.FormValue("mdn")
			if mdn == "" {
				mdn = "-"
			}
			log.Request = r.URL.RawQuery
		} else {

			reqData, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqData))
			var reqMap = make(map[string]interface{})
			json.Unmarshal(reqData, &reqMap)
			var ok bool
			mdn, ok = reqMap["msisdn"].(string)
			if !ok {
				mdn, ok = reqMap["sourcemdn"].(string)
				if !ok {
					mdn = "-"
				}
			}
			log.Request = reqMap
		}

		log.MDN = mdn

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, r)

		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		code := recorder.Code
		resData, _ := io.ReadAll(recorder.Body)

		var resMap = make(map[string]interface{})
		json.Unmarshal(resData, &resMap)
		errorcode, okErrocode := resMap["errorcode"].(int)
		if !okErrocode {
			errcodeFloat, ok := resMap["errorcode"].(float64)
			if ok {
				errorcode = int(errcodeFloat)
			} else {
				errorcodeStr, okErrocodeStr := resMap["errorcode"].(string)
				if okErrocodeStr {
					errorcode, _ = strconv.Atoi(errorcodeStr)
				}
			}
		}
		log.IsSuccessful, _ = resMap["success"].(bool)
		log.Errorcode = errorcode
		log.Message, _ = resMap["message"].(string)
		log.Response = resMap

		go callLogger(log, r)
		w.WriteHeader(code)
		w.Write(resData)
	})
}

func callURLPost(url string, req map[string]interface{}) (response map[string]interface{}) {

	response = make(map[string]interface{})
	jsonReq, _ := json.Marshal(req)
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	client := &http.Client{Timeout: 15 * time.Second, Transport: t}
	resp, err := client.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		WriteLog("Error in callURLPost " + err.Error())
		response["message"] = "Unreachable web service"
		response["errorcode"] = 500
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)
	WriteLog(string(bodyBytes))
	return
}
