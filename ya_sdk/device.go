package ya_sdk

type DeviceInfo struct {
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description,omitempty""`
	Room         string             `json:"room,omitempty`
	Type         string             `json:"type`
	Capabilities []DeviceCapability `json:"capabilities"`
	Properties   []DeviceProperty   `json:"properties"`
}

type DeviceCapability struct {
}

type DeviceProperty struct {
}
