package email

import (
	"log"
)

func SignupEmailRoutine(queue chan ActivationQItem) {
	log.Print("Starting SignupEmailRoutine")
	for {
		select {
		case data := <-queue:
			{
				SendActivationEmail(data.Owner, data.Token)
			}
		}

	}
}
