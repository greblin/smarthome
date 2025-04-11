package living_room

import "github.com/greblin/smarthome/ya_sdk"

const (
	torchereDeviceId   = "living-room-torchere"
	torchereDeviceName = "Торшер"
)

type torchere struct {
}

func InitTorchere() *torchere {
	return &torchere{}
}

func (d *torchere) GetId() string {
	return torchereDeviceId
}

func (d *torchere) Discovery() ya_sdk.DeviceInfo {
	return ya_sdk.DeviceInfo{
		Id:           d.GetId(),
		Name:         torchereDeviceName,
		Room:         roomName,
		Type:         ya_sdk.DeviceTypeLight,
		Capabilities: d.getCapabilities(),
		Properties:   d.getProperties(),
	}
}

func (d *torchere) getCapabilities() []ya_sdk.CapabilityInfo {
	return []ya_sdk.CapabilityInfo{
		{
			Type:        ya_sdk.CapabilityTypeOnOff,
			Retrievable: false,
			Parameters:  ya_sdk.OnOffParameters{Split: true},
		},
		{
			Type:        ya_sdk.CapabilityTypeRange,
			Retrievable: false,
			Parameters: ya_sdk.RangeParameters{
				Instance:     ya_sdk.CapabilityInstanceBrightness,
				RandomAccess: false,
				Range:        ya_sdk.Range{Max: 100, Min: 0, Precision: 10},
				Unit:         ya_sdk.CapabilityUnitPercent,
			},
		},
		{
			Type:        ya_sdk.CapabilityTypeColorSettings,
			Retrievable: false,
			Parameters: ya_sdk.ColorSettingsParameters{
				Temperature: &ya_sdk.ColorTemperatureRange{Max: 6500, Min: 2700},
			},
		},
	}
}

func (d *torchere) getProperties() []ya_sdk.PropertyInfo {
	return []ya_sdk.PropertyInfo{}
}

func (d *torchere) Query() ya_sdk.DeviceState {
	return ya_sdk.DeviceState{
		Id:           d.GetId(),
		Capabilities: []ya_sdk.CapabilityState{}, //все умения этого устройства не являются Retrievable, поэтому их можно не включать в ответ
		Properties:   []ya_sdk.PropertyState{},
		ErrorCode:    "",
		ErrorMessage: "",
	}
}

func (d *torchere) Actions([]ya_sdk.CapabilityState) ya_sdk.DeviceActionResult {
	return ya_sdk.DeviceActionResult{
		Id: d.GetId(),
		ActionResult: ya_sdk.ActionResult{
			Status:       "DONE",
			ErrorCode:    "",
			ErrorMessage: "",
		},
	}
}
