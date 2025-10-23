package view

type Message struct {
	Type      string      `json:"type"`
	SDP       string      `json:"sdp,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
}
