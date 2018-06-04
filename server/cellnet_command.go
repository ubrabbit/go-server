package server

import (
	"github.com/davyxu/cellnet"
	//"github.com/ubrabbit/go-server/protocol"
)

func onClientCommand(ev cellnet.Event) {
	switch msg := ev.Message().(type) {
	}
}

func onServerCommand(ev cellnet.Event) {
	switch msg := ev.Message().(type) {
	}
}
