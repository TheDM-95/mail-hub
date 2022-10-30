package mail

import (
	"time"

	"github.com/TheDM-95/mail-hub/pkg/mail/services"
)

type MailerOption func(mailer *Mailer)

func WithTimeout(timeout time.Duration) MailerOption {
	return func(m *Mailer) {
		m.timeout = timeout
	}
}

func WithMailService(ms services.MailService) MailerOption {
	return func(m *Mailer) {
		m.service = ms
	}
}
