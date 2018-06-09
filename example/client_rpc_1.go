package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

import (
	"github.com/ubrabbit/go-server/proto"
	. "github.com/ubrabbit/go-server/server"
)

var Address = "127.0.0.1:3832"

type ConnectCmd struct {
}

func ReadConsole(callback func(string)) {
	for {
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)
		callback(text)
	}
}

func (self *ConnectCmd) OnProtoCommand(c *Connect, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatACK:
		fmt.Println("custom command: ", msg)
	default:
		fmt.Println("invalid command: ", reflect.TypeOf(msg).Elem())
	}
}

func (self *ConnectCmd) OnEventTrigger(c *Connect, name string, args ...interface{}) {
	fmt.Println("CustomEventCommand:  ", c, name, args)
}

func RpcCallBack(c *Connect, msg interface{}, err error) {
	fmt.Println("RpcCallBack:   ", msg, err)
}

func main() {
	fmt.Println("start client:")

	InitConnectPool()
	handle := &ConnectCmd{}
	obj := NewTcpConnect("client", Address, handle)

	// 阻塞的从命令行获取聊天输入
	ReadConsole(func(str string) {
		fmt.Println("send: ", str)
		obj.RpcCallSync(&proto.TestChatREQ{
			Content: str,
		}, RpcCallBack)
	})

}
