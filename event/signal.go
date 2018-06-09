package event

import (
	"fmt"
)

const (
	SIGNAL_TYPE_REMOVE  = 1
	SIGNAL_TYPE_EXECUTE = 2
)

type EventSignal struct {
	Name string
	Type int
	Args []interface{}
}

func (self *EventSignal) String() string {
	return fmt.Sprintf("[Signal][%s]-%d", self.Name, self.Type)
}
