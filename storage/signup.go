package storage

import (
	"errors"
	"github.com/r3boot/go-ipam/models"
	"github.com/r3boot/go-ipam/storage/email"
	"log"
)

func RunSignup(queue chan email.ActivationQItem, owner models.Owner) (token string, err error) {
	var (
		signupRequest email.ActivationQItem
	)

	if backend.HasOwner(owner) {
		err = errors.New("RunSignup: Owner already exists: " + owner.Fullname)
		log.Print(err)
		return
	}

	if err = backend.AddOwner(owner); err != nil {
		err = errors.New("RunSignup: Failed to add owner: " + err.Error())
		log.Print(err)
		return
	}

	token = backend.GenerateToken()
	if err = backend.NewActivation(owner, token); err != nil {
		backend.DeleteOwner(owner.Username)
		err = errors.New("RunSignup: Failed to add activation token: " + err.Error())
		log.Print(err)
		token = ""
		return
	}

	signupRequest = email.ActivationQItem{token, owner}
	queue <- signupRequest

	return
}
