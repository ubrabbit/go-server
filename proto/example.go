package proto

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册
	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

type TestChatREQ struct {
	Content string
}

type TestChatACK struct {
	Content string
	Id      int64
}

// 用于消息日志打印消息内容
func (self *TestChatREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *TestChatACK) String() string { return fmt.Sprintf("%+v", *self) }

// 引用消息时，自动注册消息，这个文件可以由代码生成自动生成
func init() {

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestChatREQ)(nil)).Elem(),
		ID:    int(util.StringHash("proto.TestChatREQ")),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestChatACK)(nil)).Elem(),
		ID:    int(util.StringHash("proto.TestChatACK")),
	})
}
