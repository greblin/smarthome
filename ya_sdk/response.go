package ya_sdk

type Response struct {
	RequestId  string            `json:"request_id"`
	DiscoveryP *DiscoveryPayload `json:"payload,omitempty"`
	//  QueryP *QueryPayload `json:"payload,omitempty"`
}

type DiscoveryPayload struct {
	UserId  string       `json:"user_id"`
	Devices []DeviceInfo `json:"devices"`
}

type QueryPayload struct {
}
