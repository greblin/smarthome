package tuya

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

const (
	endpointReceiveToken = "/v1.0/token?grant_type=1"
	endpointGetScenes    = "/v2.0/cloud/scene/rule?space_id=%s"
	endpointTriggerScene = "/v2.0/cloud/scene/rule/%s/actions/trigger"
)

var (
	ErrRequestFailed               = errors.New("request failed")
	ErrResponseUnmarshalilngFailed = errors.New("response unmarshalling failed")
	ErrResultIsNotSuccess          = errors.New("result is not success")
)

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

// @see https://developer.tuya.com/en/docs/cloud/6c1636a9bd?id=Ka7kjumkoa53v
func (c *TuyaClient) receiveToken() error {
	resp, err := c.performRequest(http.MethodGet, endpointReceiveToken, nil)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, ErrRequestFailed.Error())
	}
	log.Println("tuya token response: ", string(resp))
	tokenResp := tuyaTokenResponse{}
	if err := json.Unmarshal(resp, &tokenResp); err != nil {
		log.Println(err)
		return errors.Wrap(err, ErrResponseUnmarshalilngFailed.Error())
	}
	if v := tokenResp.Result.AccessToken; v != "" {
		c.token = v
	} else {
		return ErrResultIsNotSuccess
	}
	return nil
}

// @see https://developer.tuya.com/en/docs/cloud/d7785d8964?id=Kcp2l4i0bo315
func (c *TuyaClient) GetScenes() ([]tuyaScene, error) {
	if c.token == "" {
		if err := c.receiveToken(); err != nil {
			return nil, err
		}
	}
	resp, err := c.performRequest(http.MethodGet, fmt.Sprintf(endpointGetScenes, c.spaceId), nil)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, ErrRequestFailed.Error())
	}
	log.Println("tuya getScenes response: ", string(resp))
	scenesResp := tuyaScenesResponse{}
	if err := json.Unmarshal(resp, &scenesResp); err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, ErrResponseUnmarshalilngFailed.Error())
	}
	if scenesResp.Success {
		return scenesResp.Result.List, nil
	}
	return nil, ErrResultIsNotSuccess
}

// @see https://developer.tuya.com/en/docs/cloud/89b2c8538b?id=Kcp2l54tos47r
func (c *TuyaClient) TriggerScene(sceneId string) error {
	if c.token == "" {
		if err := c.receiveToken(); err != nil {
			return err
		}
	}
	resp, err := c.performRequest(http.MethodPost, fmt.Sprintf(endpointTriggerScene, sceneId), nil)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, ErrRequestFailed.Error())
	}
	log.Println("tuya triggerScene resp:", string(resp))
	triggerResp := tuyaTriggerSceneResponse{}
	if err := json.Unmarshal(resp, &triggerResp); err != nil {
		log.Println(err)
		return errors.Wrap(err, ErrResponseUnmarshalilngFailed.Error())
	}
	if triggerResp.Success {
		return nil
	}
	return ErrResultIsNotSuccess
}
