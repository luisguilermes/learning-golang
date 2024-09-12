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
	pb.UnimplementedLaptopServiceServer
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

	// some heavy processing
	// time.Sleep(6 * time.Second)
	// if ctx.Err() == context.Canceled {
	// 	log.Print("The client has cancelled the request")
	// 	return nil, status.Error(codes.Canceled, "the client has cancelled the request")
	// }

	// if ctx.Err() == context.DeadlineExceeded {
	// 	log.Print("deadline is exceeded")
	// 	return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	// }

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

func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer,
) error {
	filter := req.GetFilter()
	log.Printf("Receive a search-laptop request with filter: %v", filter)

	err := server.Store.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}

			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("Sent laptop with ID: %s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}
