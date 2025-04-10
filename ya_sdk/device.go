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
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description,omitempty""`
	Room         string             `json:"room,omitempty"`
	Type         string             `json:"type"`
	Capabilities []DeviceCapability `json:"capabilities"`
	Properties   []DeviceProperty   `json:"properties"`
}

type DeviceCapability struct {
	Type        string `json:"type"`
	Retrievable bool   `json:"retrievable"`
	Parameters  any    `json:"parameters"`
}

type DeviceProperty struct {
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
