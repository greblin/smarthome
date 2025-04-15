package smarthome

import (
	"context"
	"fmt"
	"github.com/greblin/smarthome/devices"
	"github.com/greblin/smarthome/tuya"
	"github.com/greblin/smarthome/ya_sdk"
)

const (
	userId = "greblin" //так как навык приватный и не предполагает наличия внутренней системы пользователей, просто хардкод
)

type device interface {
	GetId() string
	Discovery() ya_sdk.DeviceInfo
	Query() ya_sdk.DeviceState
	Actions([]ya_sdk.CapabilityAction) ya_sdk.DeviceActionResult
}

type Smarthome struct {
	registry map[string]device
}

func NewSmarthome(tuyaClient *tuya.TuyaClient) *Smarthome {
	registry := make(map[string]device, 0)
	livingRoomTorchere := devices.NewTorchere(tuyaClient)
	registry[livingRoomTorchere.GetId()] = livingRoomTorchere
	return &Smarthome{
		registry: registry,
	}
}

func (m *Smarthome) ProcessRequest(ctx context.Context, request ya_sdk.Request) (*ya_sdk.Response, error) {
	rsp := &ya_sdk.Response{RequestId: request.Headers.RequestId}
	switch request.RequestType {
	case ya_sdk.RequestTypeUnlink:
		rsp.Payload = nil
		return rsp, nil
	case ya_sdk.RequestTypeDiscovery:
		rsp.Payload = m.discovery(ctx)
		return rsp, nil
	case ya_sdk.RequestTypeQuery:
		rsp.Payload = m.query(ctx, request.Payload.(ya_sdk.QueryRequestPayload))
		return rsp, nil
	case ya_sdk.RequestTypeAction:
		rsp.Payload = m.action(ctx, request.Payload.(ya_sdk.ActionRequestPayload))
		return rsp, nil
	}
	return nil, fmt.Errorf("unsupported request type: %s", request.RequestType)
}

func (m *Smarthome) discovery(ctx context.Context) *ya_sdk.DiscoveryResponsePayload {
	infoList := make([]ya_sdk.DeviceInfo, 0, len(m.registry))
	for _, d := range m.registry {
		infoList = append(infoList, d.Discovery())
	}
	return &ya_sdk.DiscoveryResponsePayload{
		UserId:  userId,
		Devices: infoList,
	}
}

func (m *Smarthome) query(ctx context.Context, payload ya_sdk.QueryRequestPayload) *ya_sdk.QueryResponsePayload {
	stateList := make([]ya_sdk.DeviceState, 0, len(payload.Devices))
	for _, d := range payload.Devices {
		if device, exists := m.registry[d.Id]; exists {
			stateList = append(stateList, device.Query())
		} else {
			stateList = append(stateList, ya_sdk.DeviceState{
				Id:           d.Id,
				ErrorCode:    "DEVICE_NOT_FOUND",
				ErrorMessage: "Данное устройство вам не принадлежит.",
			})
		}
	}
	return &ya_sdk.QueryResponsePayload{Devices: stateList}
}

func (m *Smarthome) action(ctx context.Context, payload ya_sdk.ActionRequestPayload) *ya_sdk.ActionResponsePayload {
	resultList := make([]ya_sdk.DeviceActionResult, 0, len(payload.Devices))
	for _, d := range payload.Devices {
		if device, exists := m.registry[d.Id]; exists {
			resultList = append(resultList, device.Actions(d.Capabilities))
		} else {
			resultList = append(resultList, ya_sdk.DeviceActionResult{
				Id: d.Id,
				ActionResult: ya_sdk.ActionResult{
					Status:       "ERROR",
					ErrorCode:    "DEVICE_NOT_FOUND",
					ErrorMessage: "Данное устройство вам не принадлежит.",
				},
			})
		}
	}
	return &ya_sdk.ActionResponsePayload{Devices: resultList}
}
