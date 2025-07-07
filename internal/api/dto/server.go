package dto

type RegisterServerRequest struct {
	Token string `json:"token"`
}

type SetTunnelRequest struct {
	Token     string `json:"token"`
	BridgeURL string `json:"bridge_url"`
}

type UnlinkDevicesRequest struct {
	Token string `json:"token"`
}