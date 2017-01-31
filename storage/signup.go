package storage

import (
	"errors"
	"github.com/r3boot/go-ipam/models"
	"log"
)

func RunSignup(owner models.Owner) error {
	var (
		token string
		err   error
	)

	if backend.HasOwner(owner) {
		err = errors.New("RunSignup: Owner already exists: " + *owner.Fullname)
		log.Print(err)
		return err
	}

	if err = backend.AddOwner(owner); err != nil {
		err = errors.New("RunSignup: Failed to add owner: " + err.Error())
		log.Print(err)
		return err
	}

	token = backend.GenerateToken()
	if err = backend.NewActivation(owner, token); err != nil {
		backend.DeleteOwner(owner.Username)
		err = errors.New("RunSignup: Failed to add activation token: " + err.Error())
		log.Print(err)
		return err
	}

	err = backend.SendActivationEmail(owner, token)
	if err != nil {
		backend.DeleteOwner(owner.Username)
		backend.DeleteActivation(token)
		err = errors.New("RunSignup: " + err.Error())
		log.Print(err)
		return err
	}

	return nil
}
