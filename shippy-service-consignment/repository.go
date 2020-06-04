package main

import (
	"context"
	pb "github.com/enixdark/sample/shippy-service-consignment/proto/consignment"
	"go.mongodb.org/mongo-driver/mongo"
)

type Consignment struct {
	ID string `json:"id"`
	Weight int32 `json:"weight"`
	Description string `json:"description"`
	Containers containers `json:"containers"`
	VesselID string `json:"vession_id"`
}

type Container struct {
	ID string `json:"id"`
	CustomerID string `json:"customer_id"`
	UserID string `json:"user_id"`
}

type Containers []*Container

func MarshalContainerCollection(containers []*pb.Container) []*Container {
	collection := make([]*pb.Container, 0)
	for _, container := range containers {
		collection = append(collection, UnmarshalContainer(container))
	}
	return collection
}

func UnmarshalConsignmentCollection(consignments []*Consignment) []*pb.Consignment {
	collection := make([]*pb.Consignment, 0)
	for _, consignment := range consignments {
		collection = append(collection, UnmarshalConsignment(consignment))
	}
	return collection
}

func UnmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id: container.ID,
		CustomerId: container.CustomerID,
		UserId: container.UserID,
	}
}

func MarshalContainer(container *pb.Container) *Container {
	return &Container{
		ID:         container.Id,
		CustomerID: container.CustomerId,
		UserID:     container.UserId,
	}
}

func MarshalConsignment(consignment *pb.Consignment) *Consignment {
	containers := MarshalContainerCollection(consignment.Containers)
	return &Consignment{
		ID: consignment.Id,
		Weight: consignment.Weight,
		Description: consignment.Description,
		Containers: containers,
		VesselID: consignment.VesselId,
	}
}

func UnmarshalConsignment(consignment *Consignment) *pb.Consignment {
	return &pb.Consignment{
		Id: consignment.ID,
		Weight: consignment.Weight,
		Description: consignment.Description,
		Containers: UnmarshalContainerCollection(consignment.Containers),
		VesselId: consignment.VesselID,
	}
}

type repository interface {
	Create(ctx context.Context, consignment *Consignment) error
	GetAll(ctx context.Context) ([]*Consignment, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func (repository *MongoRepository) Create(ctx context.Context, consignment *Consignment) error {
	_, error := repository.collection.InsertOne(ctx, consignment)
	return err
}

func (repository *MongoRepository) GetAll(ctx context.Context) ([]*Consignment, error) {
	
}