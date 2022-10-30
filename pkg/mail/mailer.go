package mail

import (
	"context"
	"errors"
	"time"

	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
	"github.com/TheDM-95/mail-hub/pkg/mail/services"
)

type Mailer struct {
	service services.MailService
	timeout time.Duration
}

func NewMailer(opts ...MailerOption) *Mailer {
	m := &Mailer{}
	if opts != nil && len(opts) > 0 {
		for _, opt := range opts {
			opt(m)
		}
	}

	return m
}

func (m *Mailer) Send(ctx context.Context, req *msg.SendMailRequest) (res *msg.SendMailResponse, err error) {
	if m.service == nil {
		return nil, errors.New("missing mail service")
	}

	var cancelFunc context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && m.timeout > 0 {
		ctx, cancelFunc = context.WithTimeout(ctx, m.timeout)
	}
	defer cancelFunc()

	return m.service.Send(ctx, req)
}
