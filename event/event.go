package event

import (
	"fmt"
	"sync"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

var (
	g_EventAdd               = 1
	g_EventRemove            = 2
	g_EventCmd               = 3
	g_SignalLen              = 1000
	g_GlobalEvent *EventPool = nil
)

type EventPool struct {
	sync.Mutex
	signalQueue chan *EventSignal
	listenQueue map[string][]int
	eventPool   map[int]interface{}
}

func NewEventPool() *EventPool {
	obj := new(EventPool)
	obj.signalQueue = make(chan *EventSignal, g_SignalLen)
	obj.listenQueue = make(map[string][]int, 0)
	obj.eventPool = make(map[int]interface{}, 0)
	go obj.eventListener()
	return obj
}

type EventFunc interface {
	Name() string
	ID() int
	Execute(...interface{})
}

type EventSignal struct {
	Type int
	Name string
	Args []interface{}
}

func (self *EventSignal) String() string {
	return fmt.Sprintf("[Signal][%s]-%d", self.Name, self.Type)
}

func (self *EventPool) pushEvent(name string, t int, args ...interface{}) bool {
	self.signalQueue <- &EventSignal{Name: name, Type: t, Args: args}
	return true
}

func (self *EventPool) eventListener() {
	for {
		obj := <-self.signalQueue
		if obj == nil {
			break
		}
		defer func() {
			err := recover()
			if err != nil {
				LogError("Listen Signal Error: ", obj, err)
			}
		}()

		switch obj.Type {
		case g_EventAdd:
			event := obj.Args[0]
			self.addEvent(event)
		case g_EventRemove:
			id := obj.Args[0].(int)
			self.removeEvent(id)
		case g_EventCmd:
			self.triggerEvent(obj.Name, obj.Args...)
		default:
			LogError("Invalid Event Type: ", obj.Type)
		}
	}
}

func (self *EventPool) addEvent(obj interface{}) bool {
	self.Lock()
	defer self.Unlock()

	id := obj.(EventFunc).ID()
	name := obj.(EventFunc).Name()
	_, exists := self.eventPool[id]
	if !exists {
		self.listenQueue[name] = append(self.listenQueue[name], id)
		self.eventPool[id] = obj
		return true
	}
	return false
}

func (self *EventPool) removeEvent(id int) bool {
	self.Lock()
	defer self.Unlock()

	obj, exists := self.eventPool[id]
	if exists {
		name := obj.(EventFunc).Name()
		idx := -1
		for i, v := range self.listenQueue[name] {
			if v == id {
				idx = i
				break
			}
		}
		if idx != -1 {
			self.listenQueue[name] = append(self.listenQueue[name][:idx], self.listenQueue[name][idx+1:]...)
		}
		delete(self.eventPool, id)
	}
	return true
}

func (self *EventPool) triggerEvent(name string, args ...interface{}) {
	self.Lock()
	defer self.Unlock()
	list, exists := self.listenQueue[name]
	if !exists {
		LogError("Event Not Exists:  ", name)
		return
	}
	for _, id := range list {
		event := self.eventPool[id]
		executeEvent(event, args...)
	}
}

func executeEvent(event interface{}, args ...interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("executeEvent Error:  ", event.(EventFunc).Name(), err)
		}
	}()
	event.(EventFunc).Execute(args...)
}

func getEventPool() *EventPool {
	return g_GlobalEvent
}

func AddEvent(obj interface{}) bool {
	return getEventPool().pushEvent("AddEvent", g_EventAdd, obj)
}

func RemoveEvent(id int) bool {
	return getEventPool().pushEvent("RemoveEvent", g_EventRemove, id)
}

func TriggerEvent(name string, args ...interface{}) {
	getEventPool().pushEvent(name, g_EventCmd, args...)
}

func InitEvent() {
	fmt.Println("InitEvent")
	g_GlobalEvent = NewEventPool()
}
