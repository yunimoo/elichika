package request

type FetchBootstrapRequest struct {
	BootstrapFetchTypes []int  `json:"bootstrap_fetch_types"`
	DeviceToken         string `json:"device_token"`
	DeviceName          string `json:"device_name"`
}
