package email

import (
	"github.com/r3boot/go-ipam/models"
)

type Config struct {
	Smarthost   string
	Sender      string
	SenderName  string
	NetworkName string
	Template    string
}

type TemplateData struct {
	Recipient   string
	Fullname    string
	NetworkName string
	Token       string
	Sender      string
	SenderName  string
}

type ActivationQItem struct {
	Token string
	models.Owner
}
