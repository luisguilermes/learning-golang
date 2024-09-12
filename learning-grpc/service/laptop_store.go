package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/luisguilermes/learning-golang/learning-grpc/pb"
)

var ErrAlreadyExists = errors.New("record already exists")

// LaptopStore is an interface to store laptop
type LaptopStore interface {
	// Save saves the laptop to the store
	Save(laptop *pb.Laptop) error
	// Find finds a laptop by ID
	Find(id string) (*pb.Laptop, error)
	// Search searches for laptops with filter, returns one by one via the found function
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

// InMemoryLaptopStore is a store for laptops in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save saves the laptop to the store
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// Deep copy
	other, err := deepCopy(laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %w", err)
	}

	store.data[laptop.Id] = other
	return nil
}

// Find finds a laptop by ID
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}

// Search searches for laptops with filter, returns one by one via the found function
func (store *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *pb.Filter,
	found func(laptop *pb.Laptop) error,
) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		// time.Sleep(time.Second) // Simulate slow search
		// log.Print("checking laptop id: ", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return errors.New("search cancelled")
		}

		if isQualified(filter, laptop) {
			// Deep copy
			other, err := deepCopy(laptop)
			if err != nil {
				return fmt.Errorf("cannot copy laptop data: %w", err)
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCore() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinMemory()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3 // 8 = 2ˆ3
	case pb.Memory_KILOBYTE:
		return value << 13 // 1024 * 8 = 2ˆ10 * 2ˆ3 = 2ˆ13
	case pb.Memory_MEGABYTE:
		return value << 23 // 1024 * 1024 * 8 = 2ˆ20 * 2ˆ3 = 2ˆ23
	case pb.Memory_GIGABYTE:
		return value << 33 // 1024 * 1024 * 1024 * 8 = 2ˆ30 * 2ˆ3 = 2ˆ33
	case pb.Memory_TERABYTE:
		return value << 43 // 1024 * 1024 * 1024 * 1024 * 8 = 2ˆ40 * 2ˆ3 = 2ˆ43
	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return other, nil
}
