package ya_sdk

const (
	DeviceTypeLight              = "devices.types.light"
	CapabilityTypeOnOff          = "devices.capabilities.on_off"
	CapabilityTypeRange          = "devices.capabilities.range"
	CapabilityTypeColorSettings  = "devices.capabilities.color_setting"
	CapabilityInstanceBrightness = "brightness"
	CapabilityUnitPercent        = "unit.percent"
)

type DeviceInfo struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description,omitempty""`
	Room         string           `json:"room,omitempty"`
	Type         string           `json:"type"`
	Capabilities []CapabilityInfo `json:"capabilities,omitempty"`
	Properties   []PropertyInfo   `json:"properties,omitempty"`
}

type PropertyInfo struct {
}

type PropertyState struct {
}

type CapabilityInfo struct {
	Type        string `json:"type"`
	Retrievable bool   `json:"retrievable"`
	Parameters  any    `json:"parameters"`
}

type OnOffParameters struct {
	Split bool `json:"split"`
}

type RangeParameters struct {
	Instance     string `json:"instance"`
	RandomAccess bool   `json:"random_access"`
	Range        Range  `json:"range"`
	Unit         string `json:"unit"`
}

type Range struct {
	Max       int `json:"max"`
	Min       int `json:"min"`
	Precision int `json:"precision"`
}

type ColorSettingsParameters struct {
	Temperature *ColorTemperatureRange `json:"temperature_k,omitempty"`
}

type ColorTemperatureRange struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type CapabilityState struct {
	Instance string `json:"instance"`
	Value    any    `json:"value"`
}

type DeviceState struct {
	Id           string            `json:"id"`
	Capabilities []CapabilityState `json:"capabilities"`
	Properties   []PropertyState   `json:"properties"`
	ErrorCode    string            `json:"error_code"`
	ErrorMessage string            `json:"error_message"`
}
