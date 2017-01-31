package email

import (
	"bytes"
	"errors"
	"github.com/r3boot/go-ipam/models"
	"io"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

var (
	config Config
)

func Setup(cfg Config) {
	config = cfg
}

func renderTemplate(owner models.Owner, token string) (content bytes.Buffer, err error) {
	var (
		tmpl *template.Template
		data TemplateData
	)

	tmpl, err = template.New("signup-template").Parse(SignupTemplate)
	if err != nil {
		return content, errors.New("renderTemplate: Failed to generate template: " + err.Error())
	}

	data = TemplateData{
		Recipient:   *owner.Email,
		Fullname:    *owner.Fullname,
		NetworkName: config.NetworkName,
		Sender:      config.Sender,
		SenderName:  config.SenderName,
		Token:       token,
	}

	err = tmpl.Execute(&content, data)
	if err != nil {
		return content, errors.New("renderTemplate: Failed to render template: " + err.Error())
	}

	return content, nil
}

func SendActivationEmail(owner models.Owner, token string) error {
	var (
		err        error
		thishost   string
		datawriter io.WriteCloser
		content    bytes.Buffer
		conn       *smtp.Client
	)

	if content, err = renderTemplate(owner, token); err != nil {
		err = errors.New("SendActivationEmail: " + err.Error())
		log.Print(err)
		return err
	}

	if conn, err = smtp.Dial(config.Smarthost); err != nil {
		err = errors.New("SendActivationEmail: SMTP Connect failed: " + err.Error())
		log.Print(err)
		return err
	}
	defer conn.Close()

	if thishost, err = os.Hostname(); err != nil {
		err = errors.New("SendActivationEmail: gethostname failed: " + err.Error())
		log.Print(err)
		return err
	}

	if err = conn.Hello(thishost); err != nil {
		err = errors.New("SendActivationEmail: HELO/EHLO failed: " + err.Error())
		log.Print(err)
		return err
	}

	if err = conn.Mail(config.Sender); err != nil {
		err = errors.New("SendActivationEmail: MAIL FROM failed: " + err.Error())
		log.Print(err)
		return err
	}

	if err = conn.Rcpt(*owner.Email); err != nil {
		err = errors.New("SendActivationEmail: RCTP TO failed: " + err.Error())
		log.Print(err)
		return err
	}

	if datawriter, err = conn.Data(); err != nil {
		err = errors.New("SendActivationEmail: DATA failed: " + err.Error())
		log.Print(err)
		return err
	}
	defer datawriter.Close()

	datawriter.Write(content.Bytes())

	conn.Quit()

	log.Print("SendActivationEmail: Email sent to " + *owner.Email + " via " + config.Smarthost)

	return nil
}
