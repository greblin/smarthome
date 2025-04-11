package main

import (
	"context"
	"github.com/greblin/smarthome/alice"
	"github.com/greblin/smarthome/ya_sdk"
)

func Handler(ctx context.Context, request ya_sdk.Request) (*ya_sdk.Response, error) {
	aliceModule := alice.NewAliceModule()
	return aliceModule.ProcessRequest(ctx, request)
}
