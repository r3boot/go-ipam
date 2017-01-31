package email

import (
	"bytes"
	"errors"
	"github.com/r3boot/go-ipam/models"
	"log"
	"net/smtp"
	"text/template"
)

var (
	config Config
)

func Setup(cfg Config) {
	config = cfg
}

func SendActivationEmail(owner models.Owner, token string) error {
	var (
		err        error
		tmpl       *template.Template
		data       TemplateData
		content    bytes.Buffer
		recipients []string
	)

	tmpl, err = template.New("signup-template").Parse(SignupTemplate)
	if err != nil {
		return errors.New("SendActivationEmail: Failed to generate template: " + err.Error())
	}

	data = TemplateData{
		Fullname:    *owner.Fullname,
		NetworkName: config.NetworkName,
		SenderName:  config.SenderName,
		Token:       token,
	}

	err = tmpl.Execute(&content, data)
	if err != nil {
		return errors.New("SendActivationEmail: Failed to render template: " + err.Error())
	}

	recipients = []string{*owner.Email}

	err = smtp.SendMail(config.Smarthost, nil, config.Sender, recipients, content.Bytes())
	if err != nil {
		return errors.New("SendActivationEmail: Failed to send email: " + err.Error())
	}

	log.Print("SendActivationEmail: Email sent to " + *owner.Email + " via " + config.Smarthost)

	return nil
}
