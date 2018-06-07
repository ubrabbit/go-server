package main

import (
	"time"
)

import (
	. "github.com/ubrabbit/go-server/common"
	. "github.com/ubrabbit/go-server/event"
)

type EventTest_1 struct {
	eventName string
	eventID   int
}

func (self *EventTest_1) ID() int {
	return self.eventID
}

func (self *EventTest_1) Name() string {
	return self.eventName
}

func (self *EventTest_1) Execute(args ...interface{}) {
	v1 := args[0].(int)
	v2 := args[1].(string)
	LogInfo("EventTest_1  Execute:")
	LogInfo("v1:  ", v1)
	LogInfo("v2:  ", v2)
}

func main() {
	LogInfo("Start")
	InitEvent()

	AddEvent(&EventTest_1{eventName: "test", eventID: 10086})
	TriggerEvent("test", 111, "hello")

	time.Sleep(1 * time.Millisecond)
	RemoveEvent(10086)
	TriggerEvent("test", 222, "hello22")

	time.Sleep(3 * time.Second)
	LogInfo("Finished")
}
