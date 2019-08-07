package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/huarngpa/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository simulates the use of a datastore
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()

	updated := append(repo.consignments, consignment)
	repo.consignments = updated

	repo.mu.Unlock()
	return consignment, nil
}

// service implements all of the methods that satisfy consignment.proto
type service struct {
	repo repository
}

// CreateConsignment is our server method that takes a context and request
// and creates a new consignment
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	consignment, err := s.repo.Create(req)

	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {

	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{repo})

	// Used to provide information about publicly-accessible gRPC services on a gRPC server
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
