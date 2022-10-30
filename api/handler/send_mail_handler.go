package handler

import (
	"encoding/json"
	"errors"
	"github.com/TheDM-95/mail-hub/util/constant"
	"github.com/TheDM-95/mail-hub/util/publisher"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"net/http"

	"github.com/TheDM-95/mail-hub/pkg/mail/msg"
)

type sendResponse struct {
	Success bool `json:"success"`
}

type SendMailHandler struct {
}

func NewSendMailHandler() *SendMailHandler {
	h := &SendMailHandler{}

	return h
}

func (h *SendMailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	sendReq := &msg.SendMailRequest{}

	err := json.NewDecoder(r.Body).Decode(sendReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validateSendRequest(sendReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := publisher.GetKafkaPublisher()
	sendPayload, err := json.Marshal(sendReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.Publish(constant.QueueTopicSendMail, &message.Message{UUID: watermill.NewUUID(), Payload: sendPayload})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &sendResponse{Success: true}
	payload, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)

	return
}

func (h *SendMailHandler) validateSendRequest(req *msg.SendMailRequest) error {
	if req.To == nil || req.To.Email == "" {
		return errors.New("missing recipient")
	}

	if req.Text == "" && req.Html == "" {
		return errors.New("missing content")
	}

	return nil
}
