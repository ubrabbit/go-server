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

func CustomOnCommand(c *ConnectUnit, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatACK:
		fmt.Println("custom command: ", msg)
	default:
		fmt.Println("invalid command: ", msg)
	}
}

func main() {
	fmt.Println("start client:")

	InitConnectPool()
	obj := NewTcpConnect("client", Address, CustomOnCommand)
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
