package event

var (
	g_GlobalEvent *EventPool = nil
)

type EventPool struct {
	Pool map[string]*EventPool
}

func NewEventPool() *EventPool {
	obj := new(EventPool)
	obj.Pool = make(map[string]*Event, 0)
	return obj
}

type Event struct {
	Name string
	Call func(...interface{})
	Args []interface{}
}

func newEvent(name string, f func(...interface{}), args ...interface{}) *Event {
	obj := Event{Name: name, Call: f, Args: args}
	return obj
}

func (self *EventPool) AddEvent(name string, f func(...interface{}), args ...interface{}) bool {
	_, exists := self.Pool[name]
	if !exists {
		obj := newEvent(name, f, args...)
		self.Pool[name] = obj
		return true
	}
	return false
}

func (self *EventPool) RemoveEvent(name string) bool {
	_, exists := self.Pool[name]
	if exists {
		delete(self.Pool, name)
	}
	return true
}

func (self *EventPool) TriggerEvent(name string, params ...interface{}) {

}

func InitEvent() {
	g_GlobalEvent = NewEventPool()
}
