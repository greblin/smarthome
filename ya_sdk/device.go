package ya_sdk

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
)

const (
	DeviceTypeLight               = "devices.types.light"
	CapabilityTypeOnOff           = "devices.capabilities.on_off"
	CapabilityTypeRange           = "devices.capabilities.range"
	CapabilityTypeColorSettings   = "devices.capabilities.color_setting"
	CapabilityInstanceOn          = "on"
	CapabilityInstanceBrightness  = "brightness"
	CapabilityInstanceTemperature = "temperature_k"
	CapabilityUnitPercent         = "unit.percent"
)

type DeviceInfo struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description,omitempty""`
	Room         string           `json:"room,omitempty"`
	Type         string           `json:"type"`
	Capabilities []CapabilityInfo `json:"capabilities,omitempty"`
	Properties   []PropertyInfo   `json:"properties,omitempty"`
}

type PropertyInfo struct {
}

type PropertyState struct {
}

type CapabilityInfo struct {
	Type        string `json:"type"`
	Retrievable bool   `json:"retrievable"`
	Parameters  any    `json:"parameters"`
}

type OnOffParameters struct {
	Split bool `json:"split"`
}

type RangeParameters struct {
	Instance     string `json:"instance"`
	RandomAccess bool   `json:"random_access"`
	Range        Range  `json:"range"`
	Unit         string `json:"unit"`
}

type Range struct {
	Max       int `json:"max"`
	Min       int `json:"min"`
	Precision int `json:"precision"`
}

type ColorSettingsParameters struct {
	Temperature *ColorTemperatureRange `json:"temperature_k,omitempty"`
}

type ColorTemperatureRange struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type CapabilityState struct {
	Type  string `json:"type"`
	State struct {
		Instance string `json:"instance"`
		Value    any    `json:"value"`
	} `json:"state"`
}

type DeviceState struct {
	Id           string            `json:"id"`
	Capabilities []CapabilityState `json:"capabilities"`
	Properties   []PropertyState   `json:"properties"`
	ErrorCode    string            `json:"error_code"`
	ErrorMessage string            `json:"error_message"`
}

type DeviceActionResult struct {
	Id           string       `json:"id"`
	ActionResult ActionResult `json:"action_result"`
}

type ActionResult struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func CreateDeviceActionResult(deviceId string, err error, errorCode string, errorMessage string) DeviceActionResult {
	r := DeviceActionResult{
		Id:           deviceId,
		ActionResult: ActionResult{},
	}
	if err == nil {
		r.ActionResult.Status = "DONE"
	} else {
		r.ActionResult.Status = "ERROR"
		r.ActionResult.ErrorCode = errorCode
		r.ActionResult.ErrorMessage = errorMessage
	}
	return r
}

func (s *CapabilityState) UnmarshalJSON(data []byte) error {
	sShadow := struct {
		Type  string `json:"type"`
		State struct {
			Instance string `json:"instance"`
			Value    any    `json:"value"`
		} `json:"state"`
	}{}
	if err := json.Unmarshal(data, &sShadow); err != nil {
		return errors.Wrap(err, "on unmarshal capability state shadow")
	}
	s.Type = sShadow.Type
	s.State.Instance = sShadow.State.Instance
	log.Println(string(data))
	switch s.Type {
	case CapabilityTypeOnOff:
		switch s.State.Instance {
		case CapabilityInstanceOn:
			if v, ok := sShadow.State.Value.(bool); ok {
				s.State.Value = v
				return nil
			} else {
				return errors.Errorf("bad value type for type-instance pair: %s %s", s.Type, s.State.Instance)
			}
		}
	case CapabilityTypeColorSettings:
		switch s.State.Instance {
		case CapabilityInstanceTemperature:
			if v, ok := sShadow.State.Value.(float64); ok {
				s.State.Value = int(v)
				return nil
			} else {
				return errors.Errorf("bad value type for type-instance pair: %s %s", s.Type, s.State.Instance)
			}
		}
	}

	return errors.Errorf("unsupported capability type-instance pair: %s %s", s.Type, s.State.Instance)
}
