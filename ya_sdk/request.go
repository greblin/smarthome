package ya_sdk

import (
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	RequestTypeUnlink    = "unlink"
	RequestTypeDiscovery = "discovery"
	RequestTypeQuery     = "query"
	RequestTypeAction    = "action"
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
	Devices []struct {
		Id           string            `json:"id"`
		Capabilities []CapabilityState `json:"capabilities"`
	} `json:"devices"`
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
