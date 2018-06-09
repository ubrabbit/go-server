package event

import (
	"fmt"
)

type EventSignal struct {
	Name string
	Args []interface{}
}

func (self *EventSignal) String() string {
	return fmt.Sprintf("[Signal][%s]", self.Name)
}
