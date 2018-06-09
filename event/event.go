package event

import (
	"container/list"
	"fmt"
	"sync"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

const (
	g_DefaultQueueLen = 1000 //默认队列长度
)

var (
	g_GlobalEvent *EventPool = nil
)

type EventPool struct {
	sync.Mutex
	signalQueue chan *EventSignal
	listenQueue map[string]*list.List
	eventPool   map[int]interface{}
	isClear     bool
}

func NewEventPool(queueLen int) *EventPool {
	obj := new(EventPool)
	obj.signalQueue = make(chan *EventSignal, queueLen)
	obj.listenQueue = make(map[string]*list.List, 0)
	obj.eventPool = make(map[int]interface{}, 0)
	obj.isClear = false
	go obj.eventListener()
	return obj
}

type EventFunc interface {
	Name() string
	ID() int
	Execute(...interface{})
}

func (self *EventPool) Clear() {
	self.Lock()
	defer self.Unlock()

	self.isClear = true
	self.signalQueue <- nil
}

func (self *EventPool) AddEvent(obj interface{}) (err error) {
	self.Lock()
	defer self.Unlock()

	if self.isClear {
		return newError(ERROR_EVENT_CLEARED, "EventPool Has Cleared")
	}
	id := obj.(EventFunc).ID()
	name := obj.(EventFunc).Name()
	_, exists := self.listenQueue[name]
	if !exists {
		self.listenQueue[name] = list.New()
	}
	_, exists = self.eventPool[id]
	if !exists {
		self.listenQueue[name].PushBack(id)
		self.eventPool[id] = obj
		return nil
	}
	return newError(ERROR_EVENT_ADD_EXISTS, "Event Has Add Before")
}

func (self *EventPool) RemoveEvent(id int) error {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError("RemoveEvent Error: ", id, err)
		}
		self.Unlock()
	}()

	obj, exists := self.eventPool[id]
	if exists {
		name := obj.(EventFunc).Name()
		e := self.getEvent(name, id)
		if e == nil {
			return nil
		}
		delete(self.eventPool, id)
		self.listenQueue[name].Remove(e)
	}
	return nil
}

func (self *EventPool) getEvent(name string, id int) *list.Element {
	l, ok := self.listenQueue[name]
	if !ok {
		return nil
	}
	for v := l.Front(); v != nil; v = v.Next() {
		if v.Value == id {
			return v
		}
	}
	return nil
}

func (self *EventPool) TriggerEvent(name string, args ...interface{}) {
	self.signalQueue <- &EventSignal{Name: name, Args: args}
}

func (self *EventPool) eventListener() {
	for {
		obj := <-self.signalQueue
		if obj == nil {
			LogInfo(self, "eventListener finished")
			break
		}

		self.Lock()
		self.executeEvent(obj)
		self.Unlock()
	}
}

func (self *EventPool) executeEvent(obj *EventSignal) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("executeEvent Error: ", obj, err)
		}
	}()
	name := obj.Name
	l, exists := self.listenQueue[name]
	if !exists {
		return
	}
	for v := l.Front(); v != nil; v = v.Next() {
		id := v.Value.(int)
		event := self.eventPool[id]
		event.(EventFunc).Execute(obj.Args...)
	}
}

func AddEvent(obj interface{}) error {
	return g_GlobalEvent.AddEvent(obj)
}

func RemoveEvent(id int) error {
	return g_GlobalEvent.RemoveEvent(id)
}

func TriggerEvent(name string, args ...interface{}) {
	g_GlobalEvent.TriggerEvent(name, args...)
}

func InitEvent() {
	fmt.Println("InitEvent")
	g_GlobalEvent = NewEventPool(g_DefaultQueueLen)
}
