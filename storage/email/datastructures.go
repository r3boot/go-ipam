package email

type Config struct {
	Smarthost   string
	Sender      string
	SenderName  string
	NetworkName string
	Template    string
}

type TemplateData struct {
	Fullname    string
	NetworkName string
	Token       string
	SenderName  string
}
