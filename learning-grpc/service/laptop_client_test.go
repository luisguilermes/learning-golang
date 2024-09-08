package service_test

import (
	"net"
	"testing"

	"github.com/luisguilermes/learning-golang/learning-grpc/pb"
	"github.com/luisguilermes/learning-golang/learning-grpc/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // :0 means that the OS will choose a free port
	require.NoError(t, err)

	grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}
