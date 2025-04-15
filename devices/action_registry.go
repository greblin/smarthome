package devices

import "github.com/greblin/smarthome/ya_sdk"

type actionFn func(state ya_sdk.CapabilityAction) error

type actionRegistry map[string]map[string]actionFn

func newActionRegistry() actionRegistry {
	return make(actionRegistry, 0)
}

func (r actionRegistry) add(capabilityType string, capabilityInstance string, action actionFn) {
	if _, exists := r[capabilityType]; !exists {
		r[capabilityType] = make(map[string]actionFn)
	}
	r[capabilityType][capabilityInstance] = action
}

func (r actionRegistry) get(capabilityType string, capabilityInstance string) actionFn {
	if _, typeExists := r[capabilityType]; !typeExists {
		return nil
	}
	if _, instanceExists := r[capabilityType][capabilityInstance]; !instanceExists {
		return nil
	}
	return r[capabilityType][capabilityInstance]
}
