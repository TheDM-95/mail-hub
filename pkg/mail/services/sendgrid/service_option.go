package sendgrid

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
)

type ServiceOption func(service *MailService)

func WithDefaultSender(from *msg.EmailAddress) ServiceOption {
	return func(s *MailService) {
		s.defaultSender = mail.NewEmail(from.Name, from.Email)
	}
}

func WithApiKey(apiKey string) ServiceOption {
	return func(s *MailService) {
		s.client = sendgrid.NewSendClient(apiKey)
	}
}

func WithApiClient(cli *sendgrid.Client) ServiceOption {
	return func(s *MailService) {
		s.client = cli
	}
}
