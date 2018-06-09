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
	LogInfo("EventTest_1  Execute:  ", v1, v2)
}

func main() {
	LogInfo("Start")
	InitEvent()

	AddEvent(&EventTest_1{eventName: "event_1", eventID: 1})
	AddEvent(&EventTest_1{eventName: "event_2", eventID: 2})
	TriggerEvent("event_1", 111, "event_1 trigger_0")

	err := AddEvent(&EventTest_1{eventName: "event_1", eventID: 1})
	if err != nil {
		LogError("AddEvent Error:  ", err.Error())
	}

	time.Sleep(1 * time.Millisecond)
	RemoveEvent(1)
	TriggerEvent("event_1", 111, "event_1 trigger_1")
	RemoveEvent(3)
	RemoveEvent(4)
	TriggerEvent("event_2", 222, "event_2 trigger")

	time.Sleep(3 * time.Second)
	LogInfo("Finished")
}
