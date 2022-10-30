package msg

type EmailAddress struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SendMailRequest struct {
	From     *EmailAddress     `json:"from"`
	To       *EmailAddress     `json:"to"`
	Subject  string            `json:"subject"`
	Text     string            `json:"text"`
	Html     string            `json:"html"`
	Cc       []*EmailAddress   `json:"cc"`
	Bcc      []*EmailAddress   `json:"bcc"`
	ReplyTo  *EmailAddress     `json:"reply_to"`
	MetaData map[string]string `json:"meta_data"`
}

type SendMailResponse struct {
	Success bool
	Data    string
}
