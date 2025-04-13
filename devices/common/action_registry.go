package common

import "github.com/greblin/smarthome/ya_sdk"

type Action func(state ya_sdk.CapabilityState) error

type ActionRegistry map[string]map[string]Action

func NewActionRegistry() ActionRegistry {
	return make(ActionRegistry, 0)
}

func (r ActionRegistry) Add(capabilityType string, capabilityInstance string, action Action) {
	if _, exists := r[capabilityType]; !exists {
		r[capabilityType] = make(map[string]Action)
	}
	r[capabilityType][capabilityInstance] = action
}

func (r ActionRegistry) Get(capabilityType string, capabilityInstance string) Action {
	if _, typeExists := r[capabilityType]; !typeExists {
		return nil
	}
	if _, instanceExists := r[capabilityType][capabilityInstance]; !instanceExists {
		return nil
	}
	return r[capabilityType][capabilityInstance]
}
