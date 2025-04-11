package ya_sdk

type Response struct {
	RequestId string `json:"request_id"`
	Payload   any    `json:"payload,omitempty"`
}

type DiscoveryResponsePayload struct {
	UserId  string       `json:"user_id"`
	Devices []DeviceInfo `json:"devices"`
}

type QueryResponsePayload struct {
	Devices []DeviceState `json:"devices"`
}
