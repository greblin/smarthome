package main

import (
	"context"
	"github.com/greblin/smarthome/alice"
	"github.com/greblin/smarthome/tuya"
	"github.com/greblin/smarthome/ya_sdk"
	"os"
)

var aliceModule *alice.AliceModule = nil

func Handler(ctx context.Context, request ya_sdk.Request) (*ya_sdk.Response, error) {
	if aliceModule == nil {
		tuyaClient := tuya.NewTuyaClient(os.Getenv("TUYA_HOST"), os.Getenv("TUYA_CLIENT_ID"), os.Getenv("TUYA_SECRET"), os.Getenv("TUYA_SPACE_ID"))
		aliceModule = alice.NewAliceModule(tuyaClient)
	}
	return aliceModule.ProcessRequest(ctx, request)
}
