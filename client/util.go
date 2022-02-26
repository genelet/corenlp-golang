package client

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/encoding/protowire"
)

// Unmarshal coreNLP profile data
//
func BytesUnmarshal(data []byte, msg protoreflect.ProtoMessage) error {
	bs, n := protowire.ConsumeBytes(data)
	if n < 0 {
		return protowire.ParseError(n)
	}
	return proto.Unmarshal(bs, msg)
}
