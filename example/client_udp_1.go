package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

import (
	. "github.com/ubrabbit/go-public/common"
	. "github.com/ubrabbit/go-server/cellnet"
	"github.com/ubrabbit/go-server/proto"
)

var Address = "127.0.0.1:3832"

type ConnectCmd struct {
}

func (self *ConnectCmd) OnProtoCommand(c *UdpConnect, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		LogInfo("TestChatREQ:  %v", msg.Content)
	case *proto.TestConnect:
		LogInfo("TestConnect:  %v", msg)
	default:
		LogError("Invalid Command:  %v", msg)
	}
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

func main() {
	fmt.Println("start client:")

	handle := &ConnectCmd{}
	obj := NewUdpConnect("client", Address, handle)
	obj.PacketSend(&proto.TestConnect{Hello: "aaaaaaaaa", Account: "ubrabbit2", Password: "123456"})
	// 阻塞的从命令行获取聊天输入
	ReadConsole(func(str string) {
		fmt.Println("send: ", str)
		if str == "close" {
			return
		}
		obj.PacketSend(&proto.TestChatREQ{
			Content: str,
		})
	})

}
