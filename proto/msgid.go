// Generated by github.com/davyxu/cellnet/protoc-gen-msg
// DO NOT EDIT!
// Source: example.proto

package proto

import (
	"github.com/davyxu/cellnet"
	"reflect"
	_ "github.com/davyxu/cellnet/codec/gogopb"
	"github.com/davyxu/cellnet/codec"
)

func init() {

	// example.proto
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*TestConnect)(nil)).Elem(),
		ID:    15276,
	})
}
