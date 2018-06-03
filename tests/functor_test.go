package tests

import (
	"fmt"
	"testing"

	. "github.com/ubrabbit/go-server/common"
)

func callTest(args ...interface{}) {
	v1 := args[0].(int)
	v2 := args[1].(string)
	v3 := args[2].(string)
	fmt.Printf("Functor callback: %d %s %s\n", v1, v2, v3)
}

func TestFunctor(t *testing.T) {
	fmt.Printf("\n\n=====================  TestFunctor  =====================\n")

	obj := NewFunctor("test", callTest, 1, "2")
	obj.Call("3")
}
