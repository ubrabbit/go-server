package proto

import (
	"github.com/davyxu/cellnet"

	"fmt"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

type TestWSJson struct {
	Msg   string
	Value int32
}

func (self *TestWSJson) String() string { return fmt.Sprintf("%+v", *self) }

// 将消息注册到系统
func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*TestWSJson)(nil)).Elem(),
		ID:    int(util.StringHash("main.TestWSJson")),
	})
}
