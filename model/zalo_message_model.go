package model

type MesageResponse struct {
	Data    []DataMessage `json:"data"`
	Error   int           `json:"error"`
	Message string        `json:"message"`
}
type DataMessage struct {
	Src             int    `json:"src"`
	Time            int64  `json:"time"`
	Type            string `json:"type"`
	URL             string `json:"url,omitempty"`
	MessageID       string `json:"message_id"`
	FromID          int64  `json:"from_id"`
	ToID            int64  `json:"to_id"`
	FromDisplayName string `json:"from_display_name"`
	FromAvatar      string `json:"from_avatar"`
	ToDisplayName   string `json:"to_display_name"`
	ToAvatar        string `json:"to_avatar"`
	Message         string `json:"message,omitempty"`
	Thumb           string `json:"thumb,omitempty"`
	Description     string `json:"description,omitempty"`
	Location        string `json:"location,omitempty"`
}

type SendMessageBody struct {
	Recipient SendMessageBodyUser    `json:"recipient"`
	Message   SendMessageBodyMessage `json:"message"`
}
type SendMessageBodyAttachment struct {
	Recipient SendMessageBodyUser    `json:"recipient"`
	Message   SendMessageBodyMessageAttachment `json:"message"`
}
type SendMessageWithElement struct {
	Recipient SendMessageBodyUser    `json:"recipient"`
	Message   SendMessageBodyMessage `json:"message"`
}

type Attachment struct {
	Type    string `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	TemplateType string     `json:"template_type"`
	Elements     []Elements `json:"elements"`
}

type Elements struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURL string `json:"image_url"`
	DefaultAction DefaultAction `json:"default_action"`
}

type DefaultAction struct {
	Type string `json:"type"`
	Payload string `json:"payload"`
	URL string	`json:"url"`
}

type Buttons struct {
}

type SendMessageBodyUser struct {
	UserID string `json:"user_id"`
}

type SendMessageBodyMessage struct {
	Text string `json:"text"`
}

type SendMessageBodyMessageAttachment struct {
	Text string `json:"text"`
	Attachment Attachment `json:"attachment"`
}
type SendMessageResponse struct {
	Error   int                     `json:"error"`
	Message string                  `json:"message"`
	Data    SendMessageResponseData `json:"data"`
}
type SendMessageResponseData struct {
	MessageID string `json:"message_id"`
}
