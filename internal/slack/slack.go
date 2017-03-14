package slack

import "fmt"

// Action defines the shape of a Slack Action Field
// https://api.slack.com/docs/message-buttons#action_fields
type Action struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NewAction creates a new Action
func NewAction(value int) *Action {
	return &Action{
		Name:  "beerSelector",
		Text:  "Select This Beer",
		Type:  "button",
		Value: fmt.Sprintf("%d", value),
	}
}

// Attachment defines the shape of a Slack Attachment Field
// https://api.slack.com/docs/message-buttons#attachment_fields
type Attachment struct {
	Title      string    `json:"title,omitempty"`
	Text       string    `json:"text,omitempty"`
	ImageURL   string    `json:"image_url,omitempty"`
	CallbackID string    `json:"callback_id"`
	Actions    []*Action `json:"actions,omitempty"`
}

// Message defines the shape of a Slack Message
// https://api.slack.com/docs/message-formatting
type Message struct {
	ResponseType string        `json:"response_type,omitempty"`
	Text         string        `json:"text"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
}

// NewMessage creates a new Message
func NewMessage() *Message {
	return &Message{
		ResponseType: "in_channel",
		Text:         "Your Untappd Response",
		Attachments:  make([]*Attachment, 0),
	}
}

// NewEmptyMessage creates a new ephemeral Message
func NewEmptyMessage() *Message {
	return &Message{
		ResponseType: "ephemeral",
		Text:         "No Results Found",
	}
}

// ActionPayload defines the shape of a Slack Action URL Invocation Payload
// https://api.slack.com/docs/message-buttons#action_url_invocation_payload
type ActionPayload struct {
	Actions []struct {
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
}
