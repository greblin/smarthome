package main

import (
	"context"
	"fmt"
	"github.com/greblin/smarthome/devices/living_room"
	"github.com/greblin/smarthome/ya_sdk"
)

type device interface {
	GetDeviceInfo() ya_sdk.DeviceInfo
}

func Handler(ctx context.Context, request ya_sdk.Request) (*ya_sdk.Response, error) {
	rsp := &ya_sdk.Response{RequestId: request.H.RequestId}
	switch request.RequestType {
	case ya_sdk.RequestTypeDiscovery:
		if p, err := discovery(ctx, request); err == nil {
			rsp.DiscoveryP = p
			return rsp, nil
		} else {
			return nil, err
		}
	}
	return nil, fmt.Errorf("unsupported request type: %s", request.RequestType)
}

func discovery(ctx context.Context, request ya_sdk.Request) (*ya_sdk.DiscoveryPayload, error) {
	userId := "greblin"
	devices := []device{living_room.InitTorchere()}
	info := make([]ya_sdk.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		info = append(info, d.GetDeviceInfo())
	}
	return &ya_sdk.DiscoveryPayload{
		UserId:  userId,
		Devices: info,
	}, nil
}
