package main

import (
	"context"
	pb "github.com/enixdark/sample/shippy-service-vessel/proto/vessel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/operator-framework/operator-sdk/vendor/k8s.io/kubernetes/pkg/scheduler/api"
)

type repository interface {
	FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error)
	Create(ctx context.Context, vessel *Vessel) error
}

type VesselRepository struct {
	collection *mongo.Collection
}

type Specification struct {
	Capacity int32
	MaxWeight int32
}

func MarshalSpecification(spec *pb.Specification) *Specification {
	return &Specification{
		Capacity: spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

func UnmarshalSpecification(spec *Specification) *pb.Specification {
	return &pb.Specification{
		Capacity: spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

func MarshalVessel(vessel *pb.Vessel) *Vessel {
	return &Vessel{
		ID: vessel.Id,
		Capacity: vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name: vessel.Name,
		Available: vessel.Available,
		OwnerID: vessel.OwnerId,
	}
}

func UnmarshalVessel(vessel *Vessel) *pb.Vessel {
	return &pb.Vessel{
		Id: vessel.ID,
		Capacity: vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name: vessel.Name,
		Available: vessel.Available,
		OwnerId: vessel.OwnerID,
	}
}

type Vessel struct {
	ID string
	Capacity int32
	Name string
	Available bool
	OwnerID string
	MaxWeight int32
}

func (repository *VesselRepository) FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error) {
	filter := bson.D{{
		"capacity",
		bson.D{{
			"$lte",
			spec.Capacity,
		}, {
			"$lte",
			spec.MaxWeight,
		}},
	}}
	vessel := &Vessel{}
	if err := repository.collection.FindOne(ctx, filter).Decode(vessel); err != nil {
		return nil, err
	}
	return vessel, nil
}

func (repository *VesselRepository) Create(ctx context.Context, vessel *Vessel) error {
	_, err := repository.collection.InsertOne(ctx, vessel)
	return err
}