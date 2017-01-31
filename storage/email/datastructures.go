package email

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
