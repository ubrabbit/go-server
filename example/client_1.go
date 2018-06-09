package main

import (
	"bufio"
	"fmt"
	"os"
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

func (self *ConnectCmd) OnProtoCommand(c *TcpConnect, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatACK:
		fmt.Println("custom command: ", msg)
	default:
		fmt.Println("invalid command: ", msg)
	}
}

func (self *ConnectCmd) OnEventTrigger(c *TcpConnect, name string, args ...interface{}) {
	fmt.Println("CustomEventCommand:  ", c, name, args)
}

func main() {
	fmt.Println("start client:")

	handle := &ConnectCmd{}
	obj := NewTcpConnect("client", Address, handle)
	obj.PacketSend(&proto.C2SConnect{Hello: "aaaaaaaaa", Account: "ubrabbit2", Password: "123456"})
	// 阻塞的从命令行获取聊天输入
	ReadConsole(func(str string) {
		fmt.Println("send: ", str)
		if str == "close" {
			obj.Disconnect()
			return
		}
		obj.PacketSend(&proto.TestChatREQ{
			Content: str,
		})
	})

}
