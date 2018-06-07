package event

import (
	"fmt"
	"sync"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

var (
	g_TriggerLen             = 1000
	g_GlobalEvent *EventPool = nil
)

type EventPool struct {
	sync.Mutex
	listenQueue map[string][]int
	eventPool   map[int]interface{}
}

func NewEventPool() *EventPool {
	obj := new(EventPool)
	obj.listenQueue = make(map[string][]int, 0)
	obj.eventPool = make(map[int]interface{}, 0)
	return obj
}

type EventFunc interface {
	Name() string
	ID() int
	Execute(...interface{})
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
	return getEventPool().addEvent(obj)
}

func RemoveEvent(id int) bool {
	return getEventPool().removeEvent(id)
}

func TriggerEvent(name string, args ...interface{}) {
	getEventPool().triggerEvent(name, args...)
}

func InitEvent() {
	fmt.Println("InitEvent")
	g_GlobalEvent = NewEventPool()
}
