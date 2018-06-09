package event

const (
	ERROR_SUCCESS          = iota
	ERROR_EVENT_CLEARED    //事件池已经被清除
	ERROR_EVENT_ADD_EXISTS //事件已存在
)

type EventError struct {
	Code   int
	Reason string
}

func (self *EventError) Error() string {
	return self.Reason
}

func newError(code int, reason string) *EventError {
	return &EventError{Code: code, Reason: reason}
}
