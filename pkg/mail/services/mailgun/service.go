package mailgun

import (
	"context"
	"errors"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
)

type MailService struct {
	defaultSender string
	client        mailgun.Mailgun
}

func NewMailService(opts ...ServiceOption) *MailService {
	s := &MailService{}

	if opts != nil && len(opts) > 0 {
		for _, opt := range opts {
			opt(s)
		}
	}

	return s
}

func (s *MailService) GetName() string {
	return "Mailgun"
}

func (s *MailService) Send(ctx context.Context, req *msg.SendMailRequest) (res *msg.SendMailResponse, err error) {
	res = &msg.SendMailResponse{
		Success: false,
	}

	isEmptySender := req.From == nil || req.From.Email == ""
	if s.defaultSender == "" && isEmptySender {
		return nil, errors.New("missing sender")
	}

	var from string
	if isEmptySender {
		from = s.defaultSender
	} else {
		from = s.GetAddressInfo(req.From)
	}

	to := s.GetAddressInfo(req.To)
	message := s.client.NewMessage(from, req.Subject, req.Text, to)

	if len(req.Cc) > 0 {
		for _, cc := range req.Cc {
			message.AddCC(s.GetAddressInfo(cc))
		}
	}

	if len(req.Bcc) > 0 {
		for _, bcc := range req.Bcc {
			message.AddCC(s.GetAddressInfo(bcc))
		}
	}

	if req.Html != "" {
		message.SetHtml(req.Html)
	}

	if req.ReplyTo != nil {
		message.SetReplyTo(s.GetAddressInfo(req.ReplyTo))
	}

	if len(req.MetaData) > 0 {
		for k, v := range req.MetaData {
			_ = message.AddVariable(k, v)
		}
	}

	println(fmt.Sprintf("Mailgun sent mail from %s to %s with subject %s content %s", from, to, req.Subject, req.Html))

	resp, mid, err := s.client.Send(ctx, message)

	if err != nil {
		return nil, err
	}

	res.Success = true
	res.Data = fmt.Sprintf(`{"message":"%v","id":"%v"}`, resp, mid)

	return res, nil
}

func (s *MailService) GetAddressInfo(add *msg.EmailAddress) string {
	if add != nil && add.Email != "" {
		return fmt.Sprintf("%v <%v>", add.Name, add.Email)
	}
	return ""
}
