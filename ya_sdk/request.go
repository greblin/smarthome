package ya_sdk

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
)

const (
	RequestTypeUnlink    = "unlink"
	RequestTypeDiscovery = "discovery"
	RequestTypeQuery     = "query"
	RequestTypeAction    = "action"
)

const (
	valuesDeviceIdKey        = "deviceId"
	valuesCapabilityType     = "capType"
	valuesCapabilityInstance = "capInstance"
	valuesCapabilityRelative = "capRelative"
	valuesCapabilityValue    = "capValue"
)

type Request struct {
	Headers     Headers `json:"headers"`
	RequestType string  `json:"request_type"`
	Payload     any     `json:"payload"`
	ApiVersion  float64 `json:"api_version"`
}

type Headers struct {
	RequestId     string `json:"request_id"`
	Authorization string `json:"authorization"`
}

type QueryRequestPayload struct {
	Devices []struct {
		Id string `json:"id"`
	} `json:"devices"`
}

type ActionRequestPayload struct {
	Devices []ActionRequestDevice `json:"devices"`
}

type ActionRequestDevice struct {
	Id           string             `json:"id"`
	Capabilities []CapabilityAction `json:"capabilities"`
}

func CreateActionRequestFromValues(ctx context.Context, values url.Values) (Request, error) {
	deviceId := values.Get(valuesDeviceIdKey)
	if deviceId == "" {
		return Request{}, errors.New("undefined device id")
	}
	if action, err := createCapabilityActionFromValues(values); err == nil {
		return Request{
			Headers:     Headers{},
			RequestType: RequestTypeAction,
			Payload: ActionRequestPayload{
				Devices: []ActionRequestDevice{
					{
						Id:           deviceId,
						Capabilities: []CapabilityAction{action},
					},
				},
			},
		}, nil
	} else {
		return Request{}, err
	}
}

func (rq *Request) UnmarshalJSON(data []byte) error {
	rqShadow := struct {
		Headers     Headers `json:"headers"`
		RequestType string  `json:"request_type"`
		ApiVersion  float64 `json:"api_version"`
	}{}
	if err := json.Unmarshal(data, &rqShadow); err != nil {
		return errors.Wrap(err, "on unmarshal request shadow")
	}
	rq.Headers = rqShadow.Headers
	rq.RequestType = rqShadow.RequestType
	rq.ApiVersion = rqShadow.ApiVersion
	switch rqShadow.RequestType {
	case RequestTypeQuery:
		queryRqShadow := struct {
			Payload QueryRequestPayload `json:"payload"`
		}{}
		if err := json.Unmarshal(data, &queryRqShadow); err != nil {
			return errors.Wrap(err, "on unmarshal query request shadow")
		}
		rq.Payload = queryRqShadow.Payload
	case RequestTypeAction:
		actionRqShadow := struct {
			Payload ActionRequestPayload `json:"payload"`
		}{}
		if err := json.Unmarshal(data, &actionRqShadow); err != nil {
			return errors.Wrap(err, "on unmarshal action request shadow")
		}
		rq.Payload = actionRqShadow.Payload
	}
	return nil
}
