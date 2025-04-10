package living_room

import "github.com/greblin/smarthome/ya_sdk"

type torchere struct {
}

func InitTorchere() *torchere {
	return &torchere{}
}

func (d *torchere) GetDeviceInfo() ya_sdk.DeviceInfo {
	return ya_sdk.DeviceInfo{
		Id:           "living-room-torchere",
		Name:         "Торшер",
		Room:         "Гостиная",
		Type:         "devices.types.light",
		Capabilities: []ya_sdk.DeviceCapability{},
		Properties:   []ya_sdk.DeviceProperty{},
	}
}
