package serializer

import (
	"testing"

	"github.com/luisguilermes/learning-golang/learning-grpc/pb"
	"github.com/luisguilermes/learning-golang/learning-grpc/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binariFile := "laptop.bin.tmp"
	jsonFile := "laptop.json.tmp"

	laptop := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop, binariFile)
	require.NoError(t, err)

	otherLaptop := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binariFile, otherLaptop)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop, otherLaptop))

	err = WriteProtobufToJSONFile(laptop, jsonFile)
	require.NoError(t, err)
}
