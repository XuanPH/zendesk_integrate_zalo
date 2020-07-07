package model

type Metadata struct {
	Name               string   `json:"name"`
	Token              string   `json:"token"`
	Tags               []string `json:"tags"`
	Priority           string   `json:"priority"`
	Type               string   `json:"type"`
	Author             int64    `json:"author"`
	AuthorProfile      Data     `json:"author_profile"`
	ReturnUrl          string   `json:"return_url"`
	InstancePushId     string   `json:"instance_push_id" form:"instance_push_id"`
	Locale             string   `json:"locale" form:"locale"`
	Subdomain          string   `json:"subdomain" form:"subdomain"`
	ZendeskAccessToken string   `json:"zendesk_access_token" form:"zendesk_access_token"`
}

type Manifest struct {
	Name             string       `json:"name"`
	ID               string       `json:"id"`
	Author           string       `json:"author"`
	Version          string       `json:"version"`
	ChannelbackFiles bool         `json:"channelback_files"`
	PushClientId     string       `json:"push_client_id"`
	Urls             ManifestUrls `json:"urls"`
}

type ManifestUrls struct {
	AdminUI         string `json:"admin_ui"`
	PullURL         string `json:"pull_url"`
	ChannelbackURL  string `json:"channelback_url"`
	ClickthroughURL string `json:"clickthrough_url"`
	HealthcheckURL  string `json:"healthcheck_url"`
}

type IntegrationRequest struct {
	Name       string   `json:"name" form:"name"`
	ReturnUrl  string   `json:"return_url" form:"return_url"`
	Metadata   Metadata `json:"metadata" form:"metadata"`
	State      State    `json:"state" form:"state"`
	ParentId   string   `json:"thread_id" form:"thread_id"`
	Message    string   `json:"message" form:"message"`
	FileUrls   []string `json:"file_urls" form:"file_urls"`
	ExternalId string   `json:"external_id" form:"external_id"`
}

type State struct {
	Offset int `json:"offset"`
}
