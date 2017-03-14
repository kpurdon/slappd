package slack

import "fmt"

// Action ...
type Action struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NewAction creats a new action item
func NewAction(value int) *Action {
	return &Action{
		Name:  "beerSelector",
		Text:  "Select This Beer",
		Type:  "button",
		Value: fmt.Sprintf("%d", value),
	}
}

// Attachment ...
type Attachment struct {
	Title      string    `json:"title,omitempty"`
	Text       string    `json:"text,omitempty"`
	ImageURL   string    `json:"image_url,omitempty"`
	CallbackID string    `json:"callback_id"`
	Actions    []*Action `json:"actions,omitempty"`
}

// Response ...
type Response struct {
	ResponseType string        `json:"response_type,omitempty"`
	Text         string        `json:"text"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
}

// NewResponse creates a new default response
func NewResponse() *Response {
	return &Response{
		ResponseType: "in_channel",
		Text:         "Your Untappd Response",
		Attachments:  make([]*Attachment, 0),
	}
}

// NewEmptyResponse creates a new response that idicates no search results were found
func NewEmptyResponse() *Response {
	return &Response{
		ResponseType: "ephemeral",
		Text:         "No Results Found",
	}
}

// ActionPayload ...
type ActionPayload struct {
	Actions []struct {
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
}
