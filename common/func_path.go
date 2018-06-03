package common

import (
	"bufio"
	"fmt"
	"os"
)

//文件内容遍历
func SeekFile(path string) (func() (string, bool), error) {
	fobj, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fobj.Close()
		return nil, err
	}
	scanner := bufio.NewScanner(fobj)

	return func() (string, bool) {
		for scanner.Scan() {
			code := scanner.Text()
			return string(code), true
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("err is ", err)
		}
		fobj.Close()
		return "", false
	}, nil
}

/*
判断文件或文件夹是否存在
如果返回的错误为nil,说明文件或文件夹存在
如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
如果返回的错误为其它类型,则不确定是否在存在
*/
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
