package slack

// Attachment ...
type Attachment struct {
	Title    string `json:"title,omitempty"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
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

// NewEmptyResultsResponse creates a new response that idicates no search results were found
func NewEmptyResultsResponse() *Response {
	return &Response{
		ResponseType: "ephemeral",
		Text:         "No Results Found",
	}
}
