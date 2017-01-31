package postgres

import (
	"errors"
	"fmt"
	"github.com/r3boot/go-ipam/models"
	"time"
)

func NewActivation(owner models.Owner, token string) error {
	var (
		err        error
		activation models.Activation
	)

	activation = models.Activation{
		Token:          &token,
		Username:       owner.Username,
		GenerationTime: time.Now().Format(time.RFC3339),
	}

	err = db.Insert(&activation)
	if err != nil {
		fmt.Println("AddActivation: Failed to insert record: " + err.Error())
		err = errors.New("NewActivation: Failed to insert record: " + err.Error())
	}

	return err
}

func DeleteActivation(token string) error {
	var (
		err        error
		activation models.Activation
	)

	_, err = db.Model(&activation).
		Column("token").
		Where("token = ?", token).
		Delete()

	return err
}

func GetActivations() models.Activations {
	var (
		err         error
		activations models.Activations
	)

	err = db.Model(&activations).Select()
	if err != nil {
		fmt.Println("GetActivations: Select failed: " + err.Error())
		return nil
	}

	return activations
}

func GetActivation(token string) models.Activation {
	var (
		err        error
		activation models.Activation
	)

	err = db.Model(&activation).
		Where("token = ?", token).
		Select()

	if err != nil {
		fmt.Println("GetActivation: Select failed: " + err.Error())
		return models.Activation{}
	}

	return activation
}

func HasActivation(token string) bool {
	return GetActivation(token).Token != nil
}
