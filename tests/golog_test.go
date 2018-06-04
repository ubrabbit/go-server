package tests

import (
	"fmt"
	golog "github.com/davyxu/golog"
	"log"
)

import "testing"

var g_Logger = golog.New("test")

func TestGolog(t *testing.T) {
	fmt.Printf("\n\n=====================  TestGolog  =====================\n")

	err := golog.SetColorFile("test", "golog_test.json")
	if err != nil {
		log.Fatalf("SetColorFile error: %v", err)
	}

	g_Logger.Errorln("11111111")
	g_Logger.Warnln("2222222")
	g_Logger.Infoln("3333333")
}
