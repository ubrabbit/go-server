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
	g_DefaultQueueLen = 10000 //默认队列长度
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

func (self *EventPool) AddEvent(obj interface{}) error {
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
	defer self.Unlock()

	_, exists := self.eventPool[id]
	if exists {
		//通过推送队列来删除，保证先前已经进入队列但还没执行的事件都执行完
		self.signalQueue <- &EventSignal{Name: "Remove", Type: SIGNAL_TYPE_REMOVE, Args: []interface{}{id}}
	}
	return nil
}

func (self *EventPool) TriggerEvent(name string, args ...interface{}) {
	self.signalQueue <- &EventSignal{Name: name, Type: SIGNAL_TYPE_EXECUTE, Args: args}
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

func (self *EventPool) getEventElem(name string, id int) *list.Element {
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

func (self *EventPool) executeEvent(obj *EventSignal) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("executeEvent Error: ", obj, err)
		}
	}()
	name := obj.Name
	switch obj.Type {
	case SIGNAL_TYPE_REMOVE:
		id := obj.Args[0].(int)
		obj, exists := self.eventPool[id]
		if exists {
			delete(self.eventPool, id)
			name := obj.(EventFunc).Name()
			elem := self.getEventElem(name, id)
			if elem != nil {
				self.listenQueue[name].Remove(elem)
			}
		}
	case SIGNAL_TYPE_EXECUTE:
		lst, exists := self.listenQueue[name]
		if exists {
			for v := lst.Front(); v != nil; v = v.Next() {
				id := v.Value.(int)
				event := self.eventPool[id]
				event.(EventFunc).Execute(obj.Args...)
			}
		}
	default:
		LogError("Invalid EventSignal Error: ", obj)
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
