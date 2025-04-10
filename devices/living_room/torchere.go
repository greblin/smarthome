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
		Id:           torchereDeviceId,
		Name:         torchereDeviceName,
		Room:         roomName,
		Type:         ya_sdk.DeviceTypeLight,
		Capabilities: []ya_sdk.DeviceCapability{},
		Properties:   []ya_sdk.DeviceProperty{},
	}
}
