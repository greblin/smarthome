package tuya

type tuyaTokenResponse struct {
	Result struct {
		AccessToken  string `json:"access_token"`
		ExpireTime   int    `json:"expire_time"`
		RefreshToken string `json:"refresh_token"`
		UID          string `json:"uid"`
	} `json:"result"`
	Success bool  `json:"success"`
	T       int64 `json:"t"`
}

type tuyaScene struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type tuyaScenesResponse struct {
	Result struct {
		List []tuyaScene `json:"list"`
	}
	Success bool `json:"success"`
}

type tuyaTriggerSceneResponse struct {
	Success bool `json:"success"`
}
