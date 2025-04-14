package living_room

import (
	"fmt"
	"github.com/greblin/smarthome/devices/common"
	"github.com/greblin/smarthome/tuya"
	"github.com/greblin/smarthome/ya_sdk"
	"github.com/pkg/errors"
	"log"
	"math"
	"strings"
)

const (
	torchereDeviceId   = "torchere"
	torchereDeviceName = "Торшер"
)

type torchere struct {
	actionsRegistry      common.ActionRegistry
	tuyaClient           *tuya.TuyaClient
	smartLifeScenes      map[string]string
	supportedTemperature []int
}

func InitTorchere(tuyaClient *tuya.TuyaClient) *torchere {
	t := &torchere{
		tuyaClient:           tuyaClient,
		supportedTemperature: []int{2700, 3200, 5000, 5500, 6500},
	}
	registry := common.NewActionRegistry()
	registry.Add(ya_sdk.CapabilityTypeOnOff, ya_sdk.CapabilityInstanceOn, t.switchOnOff)
	registry.Add(ya_sdk.CapabilityTypeColorSettings, ya_sdk.CapabilityInstanceTemperature, t.setTemperature)
	t.actionsRegistry = registry
	return t
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
				Temperature: &ya_sdk.ColorTemperatureRange{Max: d.supportedTemperature[len(d.supportedTemperature)-1], Min: d.supportedTemperature[0]},
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

func (d *torchere) Actions(actions []ya_sdk.CapabilityState) ya_sdk.DeviceActionResult {
	for _, action := range actions {
		log.Println(action)
		if handler := d.actionsRegistry.Get(action.Type, action.State.Instance); handler != nil {
			if err := handler(action); err != nil {
				return ya_sdk.CreateDeviceActionResult(d.GetId(), err, "INTERNAL_ERROR", "Случилось что-то непонятное. Подождите немного и попробуйте ещё раз.")
			}
		} else {
			return ya_sdk.CreateDeviceActionResult(d.GetId(), errors.New("INVALID_ACTION"), "INVALID_ACTION", "Это устройство так не умеет. Попробуйте что-нибудь другое.")
		}
	}
	return ya_sdk.CreateDeviceActionResult(d.GetId(), nil, "", "")
}

func (d *torchere) getSceneId(sceneName string) (string, error) {
	if len(d.smartLifeScenes) == 0 {
		scenes, err := d.tuyaClient.GetScenes()
		if err != nil {
			return "", err
		}
		scenesMap := map[string]string{}
		for _, scene := range scenes {
			nameParts := strings.Split(scene.Name, "::")
			if len(nameParts) != 2 {
				continue
			}
			if nameParts[0] != torchereDeviceId {
				continue
			}
			scenesMap[nameParts[1]] = scene.Id
		}
		d.smartLifeScenes = scenesMap
	}
	if sceneId, exists := d.smartLifeScenes[sceneName]; exists {
		return sceneId, nil
	}
	return "", errors.Errorf("unknown scene: %s", sceneName)
}

func (d *torchere) switchOnOff(action ya_sdk.CapabilityState) error {
	state, ok := action.State.Value.(bool)
	if !ok {
		return errors.New("bad state value type")
	}
	sceneName := "on"
	if !state {
		sceneName = "off"
	}
	sceneId, err := d.getSceneId(sceneName)
	if err != nil {
		return err
	}
	d.tuyaClient.TriggerScene(sceneId)
	return nil
}

func (d *torchere) setTemperature(action ya_sdk.CapabilityState) error {
	value, ok := action.State.Value.(int)
	if !ok {
		return errors.New("bad state value type")
	}
	idx, minDiff := -1, 0.0
	for i, st := range d.supportedTemperature {
		diff := math.Abs(float64(st - value))
		if diff < minDiff || idx == -1 {
			idx = i
			minDiff = diff
		}
	}
	sceneName := fmt.Sprintf("temp%d", d.supportedTemperature[idx])
	sceneId, err := d.getSceneId(sceneName)
	if err != nil {
		return err
	}
	d.tuyaClient.TriggerScene(sceneId)
	return nil
}
