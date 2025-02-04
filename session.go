package midutils

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/motaz/codeutils"
)

func CheckMethodAndContentType(w http.ResponseWriter, r *http.Request, method string) (valid bool) {

	method = strings.ToUpper(method)
	valid = strings.ToUpper(r.Method) == method
	if !valid {
		SetStatusError(w, "Method must be "+method, ERR_INVALID_METHOD, http.StatusMethodNotAllowed)
	}
	if valid && (strings.ToLower(r.Header.Get("content-type")) != "application/json" && r.Method != http.MethodGet) {
		SetStatusError(w, "content-type must be application/json", 2, http.StatusNotAcceptable)
		valid = false
	}
	return
}

func CheckContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			SetStatusError(w, "content-type must be application/json", ERR_INVALID_CONTENTTYPE, http.StatusNotAcceptable)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type SessionType struct {
	UserID   int
	Token    string
	Username string
	IP       string
}

type CheckSessionType struct {
	ResponseType
	Session SessionType
}

func CheckToken(token string) (isValid bool, session CheckSessionType, err error) {

	url := GetConfigValue("suditurl", "http://localhost:9004/") + "/checksession"
	var req *http.Request
	req, err = codeutils.PrepareURLCall(url, "GET", nil)
	isValid = false
	if err == nil {
		header := make(map[string]string)
		header["content-type"] = "application/json"
		header["token"] = token
		codeutils.SetURLHeaders(req, header)
		result := codeutils.CallURL(req, 30)
		err = result.Err
		if result.Err == nil {
			err = json.Unmarshal(result.Content, &session)
			isValid = err == nil && session.Success
		}

	}
	return
}

func GetSession(r *http.Request) (token string, exist bool, session CheckSessionType, err error) {

	token = r.Header.Get("token")
	if token == "" {
		authHeaderArr := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeaderArr) == 2 {
			token = authHeaderArr[1]
		}
	}
	exist, session, err = CheckToken(token)

	return
}

func CheckNumber(w http.ResponseWriter, number string) (mdn string, valid bool) {

	if len(number) > 0 && (strings.HasPrefix(number, "0") || strings.HasPrefix(number, "+")) {
		number = number[1:]
	}
	sampleRegexp := regexp.MustCompile(`\D`)
	valid = !sampleRegexp.MatchString(number)
	if valid {
		prefix := GetConfigValue("countryprefix", "249")
		if strings.HasPrefix(number, prefix) {
			number = number[len(prefix):]
		}
		lengStr := GetConfigValue("mdnlength", "9")
		leng, _ := strconv.Atoi(lengStr)

		valid = len(number) == leng
	}
	if valid {
		mdn = number
	} else {
		SetStatusError(w, "Invalid MDN: "+number, ERR_INVALID_NUMBER, http.StatusBadRequest)

	}
	return
}

func CheckSession(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, exist, _, _ := GetSession(r)
		if exist {
			handler.ServeHTTP(w, r)
		} else {
			SetStatusError(w, "Invalid token", ERR_INVALID_TOKEN, http.StatusUnauthorized)
		}
	})
}
