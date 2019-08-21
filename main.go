package main

import (
	"context"
	"log"
	"sync"

	shippy "github.com/huarngpa/shippy-service-consignment/proto/consignment"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
)

type repository interface {
	Create(*shippy.Consignment) (*shippy.Consignment, error)
	GetAll() []*shippy.Consignment
}

// Repository simulates the use of a datastore
type Repository struct {
	mu           sync.RWMutex
	consignments []*shippy.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *shippy.Consignment) (*shippy.Consignment, error) {
	repo.mu.Lock()

	updated := append(repo.consignments, consignment)
	repo.consignments = updated

	repo.mu.Unlock()
	return consignment, nil
}

// GetAll consignments
func (repo *Repository) GetAll() []*shippy.Consignment {
	return repo.consignments
}

// service implements all of the methods that satisfy consignment.proto
type service struct {
	repo repository
}

// CreateConsignment is our server method that takes a context and request
// and creates a new consignment
func (s *service) CreateConsignment(ctx context.Context, req *shippy.Consignment, res *shippy.Response) error {

	consignment, err := s.repo.Create(req)

	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment

	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *shippy.GetRequest, res *shippy.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}

	srv := grpc.NewService(
		micro.Name("shippy.service.consignment"),
	)

	srv.Init()
	shippy.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
