package main

import (
	"sync"
	"fmt"
	pb "github.com/enixdark/sample/shippy-service-consignment/proto/consignment"
	"github.com/micro/go-micro/v2"
	"context"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	mu sync.RWMutex
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo repository
}

func(s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res = &pb.Response{
		Created: true, Consignment: consignment,
	}
	res.Created = true
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response)  error {
	consignments := s.repo.GetAll()
	res = &pb.Response{
		Consignments: consignments,
	}
	return nil
}

func main() {
	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("shippy.service.consignment"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}