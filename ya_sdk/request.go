package ya_sdk

const (
	RequestTypeUnlink    = "unlink"
	RequestTypeDiscovery = "discovery"
	RequestTypeQuery     = "query"
	RequestTypeAction    = "action"
)

type Request struct {
	H           Headers `json:"headers"`
	RequestType string  `json:"request_type"`
	ApiVersion  float64 `json:"api_version"`
}

type Headers struct {
	RequestId     string `json:"request_id"`
	Authorization string `json:"authorization"`
}
