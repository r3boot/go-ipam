package postgres

import (
	"errors"
	"fmt"
	"github.com/r3boot/go-ipam/models"
	"time"
)

func AddOwner(owner models.Owner) error {
	var (
		err error
	)

	err = db.Insert(&owner)
	if err != nil {
		fmt.Println("AddOwner: Failed to insert record: " + err.Error())
	}

	return err
}

func DeleteOwner(data interface{}) error {
	var (
		err      error
		owner    models.Owner
		username string
	)

	switch data.(type) {
	case string:
		username = data.(string)
	default:
		err = errors.New("DeleteOwner: Received a parameter with an unknown type")
		fmt.Println(err.Error())
		return err
	}

	_, err = db.Model(&owner).
		Column("username").
		Where("username = ?", username).
		Delete()

	return err
}

func GetOwners() models.Owners {
	var (
		err    error
		owners models.Owners
	)

	err = db.Model(&owners).Select()
	if err != nil {
		fmt.Println("GetOwners: Select failed: " + err.Error())
		return nil
	}

	return owners
}

func GetOwner(data interface{}) models.Owner {
	var (
		err      error
		owner    models.Owner
		username string
	)

	switch data.(type) {
	case string:
		username = data.(string)
	default:
		return models.Owner{}
	}

	err = db.Model(&owner).
		Where("username = ?", username).
		Select()

	if err != nil {
		fmt.Println("GetOwner: Select failed: " + err.Error())
		return models.Owner{}
	}

	return owner
}

func HasOwner(data interface{}) bool {
	return GetOwner(data).Username != nil
}

func UpdateOwner(owner models.Owner) error {
	var (
		err error
	)

	_, err = db.Model(&owner).
		OnConflict("(username) DO UPDATE").
		Set("fullname = ?", owner.Fullname).
		Set("email = ?", owner.Email).
		Insert()

	if err != nil {
		fmt.Println("UpdateOwner: Failed to upsert Owner: " + err.Error())
	}

	return err
}

func GetOwnerByApiToken(token string) models.Owner {
	var (
		err   error
		owner models.Owner
	)

	err = db.Model(&owner).
		Where("api_token = ?", token).
		Select()

	if err != nil {
		fmt.Println("GetOwnerByApiToken: Select failed: " + err.Error())
		return models.Owner{}
	}

	return owner
}

func GetOwnerBySessionToken(token string) models.Owner {
	var (
		err   error
		owner models.Owner
	)

	err = db.Model(&owner).
		Where("session_token = ?", token).
		Select()

	if err != nil {
		fmt.Println("GetOwnerBySessionToken: Select failed: " + err.Error())
		return models.Owner{}
	}

	return owner
}

func ActivateOwner(username string, ts string) error {
	var (
		err   error
		owner models.Owner
	)

	_, err = db.Model(&owner).
		Set("is_active = ?", true).
		Set("activation_time = ?", ts).
		Where("username = ?", username).
		Update()

	return err
}

func SetSessionToken(username string, token string) error {
	var (
		err   error
		owner models.Owner
		ts    string
	)

	ts = time.Now().Format(time.RFC3339)
	_, err = db.Model(&owner).
		Set("session_token = ?", token).
		Set("last_login = ?", ts).
		Where("username = ?", username).
		Update()

	return err
}
