package services

import (
	"context"
	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
)

type MailService interface {
	GetName() (svcName string)
	Send(ctx context.Context, req *msg.SendMailRequest) (res *msg.SendMailResponse, err error)
}
