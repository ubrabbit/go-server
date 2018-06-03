package tests

import (
	golog "github.com/davyxu/golog"
)

import "testing"

var log = golog.New("test")

func TestGolog(t *testing.T) {
	log.Errorln("111")
}
