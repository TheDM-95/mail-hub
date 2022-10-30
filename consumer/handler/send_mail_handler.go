package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"

	mailMsg "github.com/TheDM-95/mail-hub/pkg/mail/msg"

	"github.com/TheDM-95/mail-hub/pkg/mail"
)

type SendMailOption func(handler *SendMailHandler)

func WithMailer(m *mail.Mailer) SendMailOption {
	return func(h *SendMailHandler) {
		h.mailer = m
	}
}

type SendMailHandler struct {
	mailer *mail.Mailer
}

func NewSendMailHandler(opts ...SendMailOption) *SendMailHandler {
	h := &SendMailHandler{}

	if opts != nil && len(opts) > 0 {
		for _, opt := range opts {
			opt(h)
		}
	}

	return h
}

func (h *SendMailHandler) Handle(msg *message.Message) error {
	ctx := context.Background()
	logFields := logrus.Fields{"service": "mail-hub", "method": "send-mail"}

	sendReq := &mailMsg.SendMailRequest{}
	// Decode Payload Msg
	err := json.Unmarshal(msg.Payload, sendReq)
	if err != nil {
		logrus.WithFields(logFields).Error("Error Decode send request")
		return err
	}

	if err := h.validateSendRequest(sendReq); err != nil {
		logrus.WithFields(logFields).Warnf("Invalid send request. Details %v", err.Error())
		return nil
	}

	res, err := h.mailer.Send(ctx, sendReq)
	if err != nil {
		return err
	}

	if !res.Success {
		//Todo: Update mail log
		logrus.WithFields(logrus.Fields{
			"service": "mail-hub",
			"method":  "send-mail",
			"data":    res.Data,
		}).Warn("Send message failed")
		return nil
	}

	//Todo: Create mail log
	logrus.WithFields(logrus.Fields{
		"service": "mail-hub",
		"method":  "send-mail",
		"data":    res.Data,
	}).Info("Send message success")

	return nil
}

func (h *SendMailHandler) validateSendRequest(req *mailMsg.SendMailRequest) error {
	if req.To == nil || req.To.Email == "" {
		return errors.New("missing recipient")
	}
	if req.Text == "" && req.Html == "" {
		return errors.New("missing content")
	}

	return nil
}
