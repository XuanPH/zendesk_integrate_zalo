package config

const (
	DB_HOST           = ""
	DB_PORT           = 0
	DB_NAME           = ""
	DB_UID            = ""
	DB_PWD            = ""
	ZALO_END_POINT    = ""
	ZALO_APP_ID       = ""
	ZALO_CALLBACK_URL = ""
)

var (
	ACTION_MAPPING_CLICK = map[string][]string{
		"b3338f2f-1cdb-48a1-ac9d-f0fadeb66826": []string{"Tư vấn mua hàng"},
		"e80346ff-a0df-48f2-8863-82966ccaa1ac": []string{"Bảo hànhh"},
		"e5ba6402-249e-4e3d-a24d-472dc4cebe24": []string{"Khác"},
	}
	KEYWORK_FILTERING = map[string]string{
		"tivi":    "[@TAG:Tivi]",
		"tu lanh": "[@TAG:Tulanh]",
		"tulanh":  "[@TAG:Tulanh]",
	}
)
