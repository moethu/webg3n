package renderer

import "encoding/json"

type Message struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

// sendMessageToClient sends a message to the client
func (a *RenderingApp) sendMessageToClient(action string, value string) {
	m := &Message{Action: action, Value: value}
	msg_json, err := json.Marshal(m)
	if err != nil {
		a.Log().Error(err.Error())
		return
	}
	a.Log().Info("sending message: " + string(msg_json))
	a.c_imagestream <- []byte(string(msg_json))
}
