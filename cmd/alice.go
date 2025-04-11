package main

import (
	"context"
	"fmt"
	"github.com/greblin/smarthome/devices/living_room"
	"github.com/greblin/smarthome/ya_sdk"
)

const (
	userId = "greblin" //так как навык приватный и не предполагает наличия внутренней системы пользователей, просто хардкод
)

type device interface {
	GetId() string
	Discovery() ya_sdk.DeviceInfo
	Query() ya_sdk.DeviceState
}

func Handler(ctx context.Context, request ya_sdk.Request) (*ya_sdk.Response, error) {
	rsp := &ya_sdk.Response{RequestId: request.Headers.RequestId}
	switch request.RequestType {
	case ya_sdk.RequestTypeDiscovery:
		rsp.Payload = discovery(ctx)
		return rsp, nil
	case ya_sdk.RequestTypeQuery:
		rsp.Payload = query(ctx, request.Payload.(ya_sdk.QueryRequestPayload))
		return rsp, nil
	}
	return nil, fmt.Errorf("unsupported request type: %s", request.RequestType)
}

func initDevices() map[string]device {
	registry := make(map[string]device, 0)
	livingRoomTorchere := living_room.InitTorchere()
	registry[livingRoomTorchere.GetId()] = livingRoomTorchere
	return registry
}

func discovery(ctx context.Context) *ya_sdk.DiscoveryResponsePayload {
	registry := initDevices()
	infoList := make([]ya_sdk.DeviceInfo, 0, len(registry))
	for _, d := range registry {
		infoList = append(infoList, d.Discovery())
	}
	return &ya_sdk.DiscoveryResponsePayload{
		UserId:  userId,
		Devices: infoList,
	}
}

func query(ctx context.Context, payload ya_sdk.QueryRequestPayload) *ya_sdk.QueryResponsePayload {
	registry := initDevices()

	stateList := make([]ya_sdk.DeviceState, 0, len(payload.Devices))
	for _, d := range payload.Devices {
		if device, exists := registry[d.Id]; exists {
			stateList = append(stateList, device.Query())
		} else {
			stateList = append(stateList, ya_sdk.DeviceState{
				Id:           d.Id,
				ErrorCode:    "DEVICE_NOT_FOUND",
				ErrorMessage: "Данное устройство вам не принадлежит.",
			})
		}
	}
	return &ya_sdk.QueryResponsePayload{
		Devices: stateList,
	}
}
