package common

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
	"strings"
	"unsafe"
)

func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//最高效的字符串拼接
func JoinString(split string, code ...string) string {
	s_len := len(code)
	if s_len == 0 {
		return ""
	}
	if s_len == 1 {
		return code[0]
	}

	buf := bytes.Buffer{}
	for i := 0; i < s_len; i++ {
		tmp := bytes.Buffer{}
		_, err := tmp.WriteString(code[i])
		CheckFatal(err)
		if i < s_len-1 {
			_, err = tmp.WriteString(split)
			CheckFatal(err)
		}

		_, err = buf.WriteString(tmp.String())
		CheckFatal(err)

	}
	return buf.String()
}

func StripString(str string) string {
	return strings.TrimSpace(str)
}

func Byte2String(x []byte) string {
	return *(*string)(unsafe.Pointer(&x))
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	CheckFatal(err)
	return i
}

func IntToString(v int) string {
	s := strconv.Itoa(v)
	return s
}

func BytesToInt(buf []byte) int {
	data := int(binary.BigEndian.Uint32(buf))
	return data
}

func IntToBytes(n int) []byte {
	x := uint32(n)
	//创建一个内容是[]byte的slice的缓冲器
	//与bytes.NewBufferString("")等效
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
