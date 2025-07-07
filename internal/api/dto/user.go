package dto

type GetUserServerResponse struct {
	ServerToken string  `json:"server_token,omitempty"`
	BridgeURL   *string `json:"bridge_url"`
}

type LinkServerRequest struct {
	ServerToken string `json:"server_token"`
}