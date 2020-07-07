package model

type TransformChat struct {
	ExternalID       string                `json:"external_id"`
	Message          string                `json:"message"`
	HTMLMessage      string                `json:"html_message"`
	ThreadID         string                `json:"thread_id"`
	CreatedAt        string                `json:"created_at"`
	Author           TransformChatAuthor   `json:"author"`
	Fields           []TransformChatFields `json:"fields"`
	AllowChannelback bool                  `json:"allow_channelback"`
}
type TransformChatAuthor struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	ImageURL   string `json:"image_url"`
}
type TransformChatFields struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
}
