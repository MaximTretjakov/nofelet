package view

type SDPData struct {
	Type      string       `json:"type"`
	SDP       string       `json:"sdp,omitempty"`
	Candidate IceCandidate `json:"candidate,omitempty"`
}

type IceCandidate struct {
	Candidate        string `json:"candidate"`
	SdpMid           string `json:"sdpMid"`
	SdpMLineIndex    int    `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment,omitempty"`
}
