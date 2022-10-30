package mailgun

import (
	"github.com/mailgun/mailgun-go/v4"

	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
)

type ServiceOption func(service *MailService)

func WithDefaultSender(from *msg.EmailAddress) ServiceOption {
	return func(s *MailService) {
		s.defaultSender = s.GetAddressInfo(from)
	}
}

func WithDefaultApiClient(domain, apiKey string) ServiceOption {
	return func(s *MailService) {
		s.client = mailgun.NewMailgun(domain, apiKey)
	}
}

func WithApiClient(cli mailgun.Mailgun) ServiceOption {
	return func(s *MailService) {
		s.client = cli
	}
}
