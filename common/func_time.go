package common

import (
	"time"
)

func GetTimeString() string {
	var currentTime time.Time

	currentTime = time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05")
	return newFormat
}

func GetSecond() int64 {
	now := time.Now()
	return now.Unix()
}

func GetTime() int64 {
	now := time.Now()
	//把纳秒转换成毫秒
	return now.UnixNano() / 1000000
}

func CreateTimer(ms int) chan bool {
	ch := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * time.Duration(ms))
		ch <- true
	}()
	return ch
}
