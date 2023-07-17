package aws

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// EventMessage defines the structure of the SQS base event message
type EventMessage struct {
	Type             string           `json:"Type"`
	MessageID        string           `json:"MessageId"`
	TopicARN         string           `json:"TopicArn"`
	Message          *json.RawMessage `json:"Message"`
	Timestamp        string           `json:"Timestamp"`
	SignatureVersion string           `json:"SignatureVersion"`
	Signature        string           `json:"Signature"`
	SigningCertURL   string           `json:"SigningCertURL"`
	UnsubscribeURL   string           `json:"UnsubscribeURL"`
}

// ParseEventMessage parses aws types.Message body to an EventMessage
func ParseEventMessage(body string) (EventMessage, error) {
	var result EventMessage

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return result, errors.Wrapf(err, "can't unmarshal message :%s", body)
	}

	return result, nil
}
