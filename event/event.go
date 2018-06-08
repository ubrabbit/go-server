package event

import (
	"fmt"
	"sync"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

const (
	EventAdd     = 1 //添加事件
	EventRemove  = 2 //删除事件
	EventExecute = 3 //触发事件
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
	listenQueue map[string][]int
	eventPool   map[int]interface{}
}

func NewEventPool(queueLen int) *EventPool {
	obj := new(EventPool)
	obj.signalQueue = make(chan *EventSignal, queueLen)
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

func (self *EventPool) AddEvent(obj interface{}) bool {
	return self.pushEvent("AddEvent", EventAdd, obj)
}

func (self *EventPool) RemoveEvent(id int) bool {
	return self.pushEvent("RemoveEvent", EventRemove, id)
}

func (self *EventPool) TriggerEvent(name string, args ...interface{}) {
	self.pushEvent(name, EventExecute, args...)
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
		case EventAdd:
			event := obj.Args[0]
			self.add(event)
		case EventRemove:
			id := obj.Args[0].(int)
			self.remove(id)
		case EventExecute:
			self.trigger(obj.Name, obj.Args...)
		default:
			LogError("Invalid Event Type: ", obj.Type)
		}
	}
}

func (self *EventPool) add(obj interface{}) bool {
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

func (self *EventPool) remove(id int) bool {
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

func (self *EventPool) trigger(name string, args ...interface{}) {
	self.Lock()
	defer self.Unlock()

	list, exists := self.listenQueue[name]
	if !exists {
		LogError("Event Not Exists:  ", name)
		return
	}
	for _, id := range list {
		event := self.eventPool[id]
		execute(event, args...)
	}
}

func execute(event interface{}, args ...interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("executeEvent Error:  ", err)
		}
	}()
	event.(EventFunc).Execute(args...)
}

func getEventPool() *EventPool {
	return g_GlobalEvent
}

func AddEvent(obj interface{}) bool {
	return getEventPool().AddEvent(obj)
}

func RemoveEvent(id int) bool {
	return getEventPool().RemoveEvent(id)
}

func TriggerEvent(name string, args ...interface{}) {
	getEventPool().TriggerEvent(name, args...)
}

func InitEvent() {
	fmt.Println("InitEvent")
	g_GlobalEvent = NewEventPool(g_DefaultQueueLen)
}
