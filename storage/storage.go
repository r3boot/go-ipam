package storage

import (
	"github.com/r3boot/go-ipam/models"
	"github.com/r3boot/go-ipam/storage/postgres"
)

type Storage struct {
	Connect     func(interface{}) error
	AddOwner    func(models.Owner) error
	DeleteOwner func(interface{}) error
	GetOwner    func(interface{}) models.Owner
	HasOwner    func(interface{}) bool
	UpdateOwner func(models.Owner) error
	GetOwners   func() models.Owners
}

func Setup(cfg interface{}) *Storage {
	var (
		storage *Storage
	)

	storage = &Storage{
		Connect:     postgres.Connect,
		AddOwner:    postgres.AddOwner,
		DeleteOwner: postgres.DeleteOwner,
		GetOwner:    postgres.GetOwner,
		HasOwner:    postgres.HasOwner,
		UpdateOwner: postgres.UpdateOwner,
		GetOwners:   postgres.GetOwners,
	}

	storage.Connect(cfg)

	return storage
}
