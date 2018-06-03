package common

import (
	"log"
)

type FunctorObj struct {
	Name     string
	Args     []interface{}
	callfunc func(...interface{})
}

type Functor interface {
	Call(arg ...interface{})
}

func (f FunctorObj) Call(args ...interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("functor '%s' Call Error: '%s'", f.Name, err)
			log.Println("functor Error param: ", f.Name, f.Args)
			return
		}
	}()
	f.callfunc(append(f.Args, args...)...)
}

func NewFunctor(name string, f func(...interface{}), args ...interface{}) Functor {
	return FunctorObj{name, args, f}
}
