// Generated by github.com/davyxu/cellnet/protoc-gen-msg
// DO NOT EDIT!
// Source: command.proto
// Source: login.proto

package proto

import (
	"github.com/davyxu/cellnet"
	"reflect"
	_ "github.com/davyxu/cellnet/codec/gogopb"
	"github.com/davyxu/cellnet/codec"
)

func init() {

	// command.proto
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*C2SCommand)(nil)).Elem(),
		ID:    46953,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*S2CCommand)(nil)).Elem(),
		ID:    31593,
	})
	// login.proto
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*C2SConnect)(nil)).Elem(),
		ID:    45428,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*S2CConnect)(nil)).Elem(),
		ID:    30068,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*C2SLogin)(nil)).Elem(),
		ID:    34467,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*S2CLogin)(nil)).Elem(),
		ID:    19107,
	})
}
