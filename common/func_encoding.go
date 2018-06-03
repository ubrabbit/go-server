package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetGBKString(src string) string {
	code_rlt, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		fmt.Println("error: ", err)
		return src
	}
	return string(code_rlt)
}
