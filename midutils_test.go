package midutils

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	fmt.Println("Appname: ", GetAppName())
	mdn, _ := GetPhoneNumber("249122090303")
	fmt.Println(mdn)
}
