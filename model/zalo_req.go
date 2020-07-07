package model

type MessagePush struct {
	AppID  string `json:"app_id"`
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	UserIDByApp string `json:"user_id_by_app"`
	Recipient   struct {
		ID string `json:"id"`
	} `json:"recipient"`
	EventName string `json:"event_name"`
	Message   struct {
		Text        string `json:"text"`
		MsgID       string `json:"msg_id"`
		Attachments []struct {
			Payload struct {
				Thumbnail   string `json:"thumbnail"`
				Description string `json:"description"`
				URL         string `json:"url"`
				ID          string `json:"id"`
				Coordinates struct {
					Latitude  string `json:"latitude"`
					Longitude string `json:"longitude"`
				} `json:"coordinates"`
				Size     string `json:"size"`
				Name     string `json:"name"`
				Checksum string `json:"checksum"`
				Type     string `json:"type"`
			} `json:"payload"`
			Type string `json:"type"`
		} `json:"attachments"`
	} `json:"message"`
	Timestamp string `json:"timestamp"`
}
