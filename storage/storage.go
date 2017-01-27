package storage

import (
	"github.com/r3boot/go-ipam/models"
	"github.com/r3boot/go-ipam/storage/postgres"
)

type Storage struct {
	Connect      func(interface{}) error
	AddOwner     func(models.Owner) error
	DeleteOwner  func(interface{}) error
	GetOwner     func(interface{}) models.Owner
	HasOwner     func(interface{}) bool
	UpdateOwner  func(models.Owner) error
	GetOwners    func() models.Owners
	AddAsnum     func(models.Asnum) error
	DeleteAsnum  func(interface{}) error
	GetAsnum     func(interface{}) models.Asnum
	HasAsnum     func(interface{}) bool
	UpdateAsnum  func(models.Asnum) error
	GetAsnums    func() models.Asnums
	AddPrefix    func(models.Prefix) error
	DeletePrefix func(interface{}) error
	GetPrefix    func(interface{}) models.Prefix
	HasPrefix    func(interface{}) bool
	UpdatePrefix func(models.Prefix) error
	GetPrefixes  func() models.Prefixes
}

func Setup(cfg interface{}) *Storage {
	var (
		storage *Storage
	)

	storage = &Storage{
		Connect:      postgres.Connect,
		AddOwner:     postgres.AddOwner,
		DeleteOwner:  postgres.DeleteOwner,
		GetOwner:     postgres.GetOwner,
		HasOwner:     postgres.HasOwner,
		UpdateOwner:  postgres.UpdateOwner,
		GetOwners:    postgres.GetOwners,
		AddAsnum:     postgres.AddAsnum,
		DeleteAsnum:  postgres.DeleteAsnum,
		GetAsnum:     postgres.GetAsnum,
		HasAsnum:     postgres.HasAsnum,
		UpdateAsnum:  postgres.UpdateAsnum,
		GetAsnums:    postgres.GetAsnums,
		AddPrefix:    postgres.AddPrefix,
		DeletePrefix: postgres.DeletePrefix,
		GetPrefix:    postgres.GetPrefix,
		HasPrefix:    postgres.HasPrefix,
		UpdatePrefix: postgres.UpdatePrefix,
		GetPrefixes:  postgres.GetPrefixes,
	}

	storage.Connect(cfg)

	return storage
}
