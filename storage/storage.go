package storage

import (
	"github.com/r3boot/go-ipam/models"
	"github.com/r3boot/go-ipam/storage/email"
	"github.com/r3boot/go-ipam/storage/postgres"
)

type Storage struct {
	GenerateSalt           func(int) (string, error)
	GenerateHash           func(string, string) ([]byte, error)
	GenerateToken          func() string
	RunSignup              func(chan email.ActivationQItem, models.Owner) error
	DumpAsset              func(string) (models.Asset, error)
	NewActivation          func(models.Owner, string) error
	DeleteActivation       func(string) error
	GetActivations         func() models.Activations
	GetActivation          func(string) models.Activation
	HasActivation          func(string) bool
	SendActivationEmail    func(models.Owner, string) error
	Connect                func(interface{}) error
	AddOwner               func(models.Owner) error
	DeleteOwner            func(interface{}) error
	GetOwner               func(interface{}) models.Owner
	HasOwner               func(interface{}) bool
	UpdateOwner            func(models.Owner) error
	GetOwners              func() models.Owners
	GetOwnerByApiToken     func(string) models.Owner
	GetOwnerBySessionToken func(string) models.Owner
	ActivateOwner          func(string, string) error
	SetSessionToken        func(string, string) error
	AddAsnum               func(models.Asnum) error
	DeleteAsnum            func(interface{}) error
	GetAsnum               func(interface{}) models.Asnum
	HasAsnum               func(interface{}) bool
	UpdateAsnum            func(models.Asnum) error
	GetAsnums              func() models.Asnums
	AddPrefix              func(models.Prefix) error
	DeletePrefix           func(interface{}) error
	GetPrefix              func(interface{}) models.Prefix
	HasPrefix              func(interface{}) bool
	UpdatePrefix           func(models.Prefix) error
	GetPrefixes            func() models.Prefixes
}

var backend *Storage

func Setup(cfg interface{}, email_cfg interface{}) *Storage {
	var (
		storage *Storage
	)

	storage = &Storage{
		GenerateSalt:           GenerateSalt,
		GenerateHash:           GenerateHash,
		GenerateToken:          GenerateToken,
		RunSignup:              RunSignup,
		DumpAsset:              DumpAsset,
		NewActivation:          postgres.NewActivation,
		DeleteActivation:       postgres.DeleteActivation,
		GetActivations:         postgres.GetActivations,
		GetActivation:          postgres.GetActivation,
		HasActivation:          postgres.HasActivation,
		SendActivationEmail:    email.SendActivationEmail,
		Connect:                postgres.Connect,
		AddOwner:               postgres.AddOwner,
		DeleteOwner:            postgres.DeleteOwner,
		GetOwner:               postgres.GetOwner,
		HasOwner:               postgres.HasOwner,
		UpdateOwner:            postgres.UpdateOwner,
		GetOwners:              postgres.GetOwners,
		GetOwnerByApiToken:     postgres.GetOwnerByApiToken,
		GetOwnerBySessionToken: postgres.GetOwnerBySessionToken,
		ActivateOwner:          postgres.ActivateOwner,
		SetSessionToken:        postgres.SetSessionToken,
		AddAsnum:               postgres.AddAsnum,
		DeleteAsnum:            postgres.DeleteAsnum,
		GetAsnum:               postgres.GetAsnum,
		HasAsnum:               postgres.HasAsnum,
		UpdateAsnum:            postgres.UpdateAsnum,
		GetAsnums:              postgres.GetAsnums,
		AddPrefix:              postgres.AddPrefix,
		DeletePrefix:           postgres.DeletePrefix,
		GetPrefix:              postgres.GetPrefix,
		HasPrefix:              postgres.HasPrefix,
		UpdatePrefix:           postgres.UpdatePrefix,
		GetPrefixes:            postgres.GetPrefixes,
	}

	storage.Connect(cfg)
	email.Setup(email_cfg.(email.Config))

	backend = storage

	return storage
}
