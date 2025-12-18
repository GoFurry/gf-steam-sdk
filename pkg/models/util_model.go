package models

type SteamStoreTokenResponse struct {
	Data struct {
		WebapiToken string `json:"webapi_token"` // 成就唯一标识
	} `json:"data"`
	Success int `json:"success"` // 请求是否成功
}
