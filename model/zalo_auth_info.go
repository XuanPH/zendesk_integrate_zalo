package model

type ZaloAuthInfo struct {
	Id                 string `json:"id" db:"id"`
	OAID               string `json:"oaId" db:"oaId"`
	Name               string `json:"name" db:"name"`
	Metadata           string `json:"metadata" db:"metadata"`
	InstancePushId     string `json:"instancePushId" db:"instancePushId"`
	ZendeskAccessToken string `json:"zendeskAccessToken" db:"zendeskAccessToken"`
	SubDomain          string `json:"subdomain" db:"subdomain"`
	Locale             string `json:"locale" db:"locale"`
}
