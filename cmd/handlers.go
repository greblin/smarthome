package main

import (
	"context"
	"encoding/json"
	"github.com/greblin/smarthome/smarthome"
	"github.com/greblin/smarthome/tuya"
	"github.com/greblin/smarthome/ya_sdk"
	"log"
	"net/http"
	"os"
)

const (
	envTuyaHost     = "TUYA_HOST"
	envTuyaClientId = "TUYA_CLIENT_ID"
	envTuyaSecret   = "TUYA_SECRET"
	envTuyaSpaceId  = "TUYA_SPACE_ID"
)

var sh *smarthome.Smarthome = nil

func initSmarthome() {
	if sh == nil {
		tuyaClient := tuya.NewTuyaClient(os.Getenv(envTuyaHost), os.Getenv(envTuyaClientId), os.Getenv(envTuyaSecret), os.Getenv(envTuyaSpaceId))
		sh = smarthome.NewSmarthome(tuyaClient)
	}
}

// хендлер, вызываемый при обработке JSON-RPC запроса в платформе диалогов
// @see https://yandex.ru/dev/dialogs/smart-home/doc/ru/reference/resources
func DialogHandler(ctx context.Context, request ya_sdk.Request) (*ya_sdk.Response, error) {
	initSmarthome()
	r, err := sh.ProcessRequest(ctx, request)
	if err != nil {
		log.Println(err)
	}
	return r, err
}

// хендлер, вызываемый на получение запроса от навыка "Домовенок Кузя"
func KuzyaHandler(rw http.ResponseWriter, req *http.Request) {
	request, err := ya_sdk.CreateActionRequestFromValues(req.Context(), req.URL.Query())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	initSmarthome()
	if resp, err := sh.ProcessRequest(req.Context(), request); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		respJson, jsonErr := json.Marshal(resp)
		if jsonErr != nil {
			log.Println(jsonErr)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}
		rw.Write(respJson)
		rw.WriteHeader(http.StatusOK)
	}
}
