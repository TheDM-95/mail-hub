package sendgrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
)

type MailService struct {
	defaultSender *mail.Email
	client        *sendgrid.Client
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
	return "Sendgrid"
}

func (s *MailService) Send(ctx context.Context, req *msg.SendMailRequest) (result *msg.SendMailResponse, err error) {
	result = &msg.SendMailResponse{
		Success: false,
	}

	isEmptySender := req.From == nil || req.From.Email == ""
	if s.defaultSender == nil && isEmptySender {
		return nil, errors.New("missing sender")
	}

	var from *mail.Email
	if isEmptySender {
		from = s.defaultSender
	} else {
		from = mail.NewEmail(req.From.Name, req.From.Email)
	}
	to := mail.NewEmail("", req.To.Email)
	content := mail.NewContent("text/html", req.Html)
	email := mail.NewV3MailInit(from, req.Subject, to, content)

	personalization := email.Personalizations[0]
	personalization.SetHeader("X-Transport", "web")

	if len(req.Cc) > 0 {
		ccs := make([]*mail.Email, 0)
		for _, cc := range req.Cc {
			ccs = append(ccs, mail.NewEmail(cc.Name, cc.Email))
		}
		personalization.AddCCs(ccs...)
	}

	if len(req.Bcc) > 0 {
		bccs := make([]*mail.Email, 0)
		for _, bcc := range req.Bcc {
			bccs = append(bccs, mail.NewEmail(bcc.Name, bcc.Email))
		}
		personalization.AddBCCs(bccs...)
	}

	if req.ReplyTo != nil {
		replyTo := mail.NewEmail(req.ReplyTo.Name, req.ReplyTo.Email)
		email.SetReplyTo(replyTo)
	}

	if len(req.MetaData) > 0 {
		for k, v := range req.MetaData {
			personalization.SetCustomArg(k, v)
		}
	}

	println(fmt.Sprintf("Sendgrid sent mail from %s to %s with subject %s content %s", from.Address, to.Address, req.Subject, req.Html))

	response, err := s.client.Send(email)
	if err != nil {
		return result, err
	}
	resultRaw := make(map[string]interface{})
	result.Data = response.Body
	errMarshalResult := json.Unmarshal([]byte(response.Body), &result)
	if errMarshalResult == nil && resultRaw["message"] != nil {
		if res, ok := resultRaw["message"].(string); ok {
			if res == "success" {
				result.Success = true
				return result, nil
			} else {
				return result, nil
			}
		}
	}

	if response.StatusCode > 199 && response.StatusCode < 300 {
		result.Success = true
		return result, nil
	}

	return result, nil
}
