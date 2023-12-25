package midutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/motaz/codeutils"
)

func GetAppName() (appname string) {
	appname = os.Args[0]
	var seperator string
	seperator = string(os.PathSeparator)
	if strings.Contains(appname, seperator) {
		appname = appname[strings.LastIndex(appname, seperator)+1:]
	}
	return

}

func GetConfigValue(paramName string, defaultValue string) string {

	value := codeutils.GetConfigValue("config.ini", paramName)
	if value == "" {
		value = defaultValue
	}
	return value
}

func WriteLog(event string) {

	if GetConfigValue("debug", "no") == "yes" {
		fmt.Println(event)
	}
	codeutils.WriteToLog(event, GetAppName())
}

func GetMD5(text string) (hash string) {
	hash = codeutils.GetMD5(text)
	return
}
