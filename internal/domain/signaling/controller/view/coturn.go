package view

type TURNConfig struct {
	ICEServers []ICEServer `json:"iceServers"`
}

// ICEServer представляет собой одиночный сервер (STUN или TURN)
type ICEServer struct {
	URLs       interface{} `json:"urls"` // Используем interface{} для обработки как строки, так и массива строк
	Username   string      `json:"username,omitempty"`
	Credential string      `json:"credential,omitempty"`
}
