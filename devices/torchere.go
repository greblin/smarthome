package devices

import (
	"fmt"
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
	actionsRegistry      actionRegistry
	tuyaClient           *tuya.TuyaClient
	smartLifeScenes      map[string]string
	supportedTemperature []int
}

func NewTorchere(tuyaClient *tuya.TuyaClient) *torchere {
	t := &torchere{
		tuyaClient:           tuyaClient,
		supportedTemperature: []int{2700, 3200, 4000, 5000, 5500, 6500},
	}
	registry := newActionRegistry()
	registry.add(ya_sdk.CapabilityTypeOnOff, ya_sdk.CapabilityInstanceOn, t.switchOnOff)
	registry.add(ya_sdk.CapabilityTypeColorSettings, ya_sdk.CapabilityInstanceTemperature, t.setTemperature)
	registry.add(ya_sdk.CapabilityTypeRange, ya_sdk.CapabilityInstanceBrightness, t.changeBrightness)
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
		Room:         roomLivingRoom,
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
		Capabilities: []ya_sdk.CapabilityAction{}, //все умения этого устройства не являются Retrievable, поэтому их можно не включать в ответ
		Properties:   []ya_sdk.PropertyState{},
		ErrorCode:    "",
		ErrorMessage: "",
	}
}

func (d *torchere) Actions(actions []ya_sdk.CapabilityAction) ya_sdk.DeviceActionResult {
	for _, action := range actions {
		log.Println(action)
		if handler := d.actionsRegistry.get(action.Type, action.State.Instance); handler != nil {
			if err := handler(action); err != nil {
				return ya_sdk.CreateDeviceActionResult(d.GetId(), err, ya_sdk.ErrorCodeInternalError, ya_sdk.ErrorMessageInternalError)
			}
		} else {
			return ya_sdk.CreateDeviceActionResult(d.GetId(), ErrInvalidAction, ya_sdk.ErrorCodeInvalidAction, ya_sdk.ErrorMessageInvalidAction)
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

func (d *torchere) switchOnOff(action ya_sdk.CapabilityAction) error {
	state, ok := action.State.Value.(bool)
	if !ok {
		return ErrBadStateValueType
	}
	sceneName := sceneNameSwitchOn
	if !state {
		sceneName = sceneNameSwitchOff
	}
	return d.triggerScene(sceneName)
}

func (d *torchere) setTemperature(action ya_sdk.CapabilityAction) error {
	value, ok := action.State.Value.(int)
	if !ok {
		return ErrBadStateValueType
	}
	idx, minDiff := -1, 0.0
	for i, st := range d.supportedTemperature {
		diff := math.Abs(float64(st - value))
		if diff < minDiff || idx == -1 {
			idx = i
			minDiff = diff
		}
	}
	sceneName := fmt.Sprintf("%s%d", sceneNameChangeTempPrefix, d.supportedTemperature[idx])
	return d.triggerScene(sceneName)
}

func (d *torchere) changeBrightness(action ya_sdk.CapabilityAction) error {
	value, ok := action.State.Value.(int)
	if !ok {
		return ErrBadStateValueType
	}
	if !action.State.Relative {
		return ErrNotSupportedAbsValue
	}
	sceneName := sceneNameBrightnessInc
	if value < 0 {
		sceneName = sceneNameBrightnessDec
	}
	return d.triggerScene(sceneName)
}

func (d *torchere) triggerScene(name string) error {
	sceneId, err := d.getSceneId(name)
	if err != nil {
		return errors.Wrap(err, ErrSceneSystemError.Error())
	}
	if err := d.tuyaClient.TriggerScene(sceneId); err != nil {
		return errors.Wrap(err, ErrSceneSystemError.Error())
	}
	return nil
}
