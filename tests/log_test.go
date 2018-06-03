package tests

import (
	"fmt"
	. "github.com/ubrabbit/go-server/common"
	"testing"
)

func TestLogger(t *testing.T) {
	fmt.Printf("\n\n=====================  TestLogger  =====================\n")

	LogInfo("log---------info")
	LogWarning("log---------warning")
	LogError("log---------error")
	//LogFatal("fatal")
}
