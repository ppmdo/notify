package whatsapp

import (
	"context"

	"github.com/febriliankr/whatsapp-cloud-api"
	"github.com/pkg/errors"
)

// Service encapsulates the WhatsApp client along with internal state for storing contacts.
type WhatsApp struct {
	client       *whatsapp.Whatsapp
	phoneNumbers []string
}

// New returns a new instance of a WhatsApp notification service.
func New(token string, phoneID string) (*WhatsApp, error) {
	client := whatsapp.NewWhatsapp(token, phoneID)

	w := &WhatsApp{
		client:       client,
		phoneNumbers: []string{},
	}

	return w, nil
}

// AddReceivers takes WhatsApp contacts and adds them to the internal contacts list. The Send method will send
// a given message to all those contacts.
func (s *WhatsApp) AddReceivers(phoneNumbers ...string) {
	s.phoneNumbers = append(s.phoneNumbers, phoneNumbers...)
}

// Send takes a message subject and a message body and sends them to all previously set contacts.
func (s *WhatsApp) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message

	for _, phoneNumber := range s.phoneNumbers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := s.client.SendText(phoneNumber, fullMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to %s", phoneNumber)
			}
		}
	}

	return nil
}
