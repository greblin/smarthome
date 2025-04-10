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

func (d *torchere) GetDeviceInfo() ya_sdk.DeviceInfo {
	return ya_sdk.DeviceInfo{
		Id:   torchereDeviceId,
		Name: torchereDeviceName,
		Room: roomName,
		Type: ya_sdk.DeviceTypeLight,
		Capabilities: []ya_sdk.DeviceCapability{
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
		},
		Properties: []ya_sdk.DeviceProperty{},
	}
}
