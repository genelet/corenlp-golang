package client

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// BytesUnmarshal unmarshals coreNLP protobuf data
func BytesUnmarshal(data []byte, msg protoreflect.ProtoMessage) error {
	bs, n := protowire.ConsumeBytes(data)
	if n < 0 {
		return protowire.ParseError(n)
	}
	return proto.Unmarshal(bs, msg)
}
