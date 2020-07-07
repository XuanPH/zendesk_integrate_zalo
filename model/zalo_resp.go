package model

// Officall Account Response
type OAResponse struct {
	Data    Data   `json:"data"`
	Error   int    `json:"error"`
	Message string `json:"message"`
}
type Data struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Cover       string `json:"cover"`
	OaID        int64  `json:"oa_id"`
}

//Message

type Profile struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Data    struct {
		UserGender  int    `json:"user_gender"`
		UserID      int64  `json:"user_id"`
		UserIDByApp int64  `json:"user_id_by_app"`
		Avatar      string `json:"avatar"`
		Avatars     struct {
			Num120 string `json:"120"`
			Num240 string `json:"240"`
		} `json:"avatars"`
		DisplayName string `json:"display_name"`
		BirthDate   int    `json:"birth_date"`
		SharedInfo  string `json:"shared_info"`
	} `json:"data"`
}
