package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/luisguilermes/learning-golang/learning-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopServiceServer is the server that provides laptop services
type LaptopServer struct {
	Store LaptopStore
}

// NewLaptopServiceServer returns a new LaptopServiceServer
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

// CreateLapTop creates a new laptop.
//
// Parameters:
// - server: The LaptopServer instance.
//
// Returns:
// - error: An error if the laptop creation fails.
func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Print("Receive a create-laptop request with id: ", laptop.GetId())

	if len(laptop.Id) > 0 {
		// Check if the id is in UUID format
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		// Generate a new UUID
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	//save the laptop to the in-memory store
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("Laptop with ID %s saved to the store", laptop.GetId())
	res := &pb.CreateLaptopResponse{
		Id: laptop.GetId(),
	}
	return res, nil
}
