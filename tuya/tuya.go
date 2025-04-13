package tuya

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

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

type TuyaScene struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type tuyaScenesResponse struct {
	Result struct {
		List []TuyaScene `json:"list"`
	}
	Success bool `json:"success"`
}

type TuyaClient struct {
	host     string
	clientId string
	secret   string
	spaceId  string
	token    string
}

func NewTuyaClient(host, clientId, secret, spaceId string) *TuyaClient {
	return &TuyaClient{
		host:     host,
		clientId: clientId,
		secret:   secret,
		spaceId:  spaceId,
		token:    "",
	}
}

func (c *TuyaClient) receiveToken() error {
	resp, err := c.performRequest(http.MethodGet, "/v1.0/token?grant_type=1&aaa=bar", nil)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("resp:", string(resp))
	tokenResp := tuyaTokenResponse{}
	if err := json.Unmarshal(resp, &tokenResp); err != nil {
		log.Println(err)
		return err
	}
	if v := tokenResp.Result.AccessToken; v != "" {
		c.token = v
	} else {
		//todo
	}
	return nil
}

func (c *TuyaClient) GetScenes() ([]TuyaScene, error) {
	if c.token == "" {
		if err := c.receiveToken(); err != nil {
			return nil, err
		}
	}
	resp, err := c.performRequest(http.MethodGet, fmt.Sprintf("/v2.0/cloud/scene/rule?space_id=%s", c.spaceId), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("resp:", string(resp))
	scenesResp := tuyaScenesResponse{}
	if err := json.Unmarshal(resp, &scenesResp); err != nil {
		log.Println(err)
		return nil, err
	}
	if scenesResp.Success {
		return scenesResp.Result.List, nil
	}
	return nil, errors.New("some error") //todo
}

func (c *TuyaClient) TriggerScene(sceneId string) error {
	if c.token == "" {
		if err := c.receiveToken(); err != nil {
			return err
		}
	}
	resp, err := c.performRequest(http.MethodPost, fmt.Sprintf("/v2.0/cloud/scene/rule/%s/actions/trigger", sceneId), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("resp:", string(resp))
	return nil
}
