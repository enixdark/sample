package main

import (
	"os"
	"sync"
	"fmt"
	"log"
	pb "github.com/enixdark/sample/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/enixdark/sample/shippy-service-vessel/proto/vessel"
	"github.com/micro/go-micro/v2"
	"context"

	datastore "github.com/enixdark/sample/shippy-service-consignment/proto/consignment/datastore"
)

const (
	port = ":50051"
	defaultHost = "localhost:27017"
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
	vesselClient vesselProto.VesselService
}

func(s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindAvailable(
		context.Background(), &vesselProto.Specification{
			MaxWeight: req.Weight,
			Capacity: int32(len(req.Containers)),
		},
	)

	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)

	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment

	return nil
}



func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response)  error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("shippy.service.consignment"),
		micro.Version("latest"),
	)

	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := datastore.CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}

	defer client.Disconnect(context.Background())

	vesselClient := vesselProto.NewVesselService("shippy.service.vessel", srv.Client())

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}