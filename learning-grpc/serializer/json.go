package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON converts a protobuf message to JSON string
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		Indent:          "  ",
	}

	bjson, err := marshaler.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(bjson), nil
}
