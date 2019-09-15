package renderer

import "encoding/json"

// Response Message to frontend
type Message struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

// respondToClient sends a message upstream
func (a *RenderingApp) respondToClient(action string, value string) {
	m := &Message{Action: action, Value: value}
	msg_json, err := json.Marshal(m)
	if err != nil {
		a.log.Error(err.Error())
		return
	}
	a.log.Info("sending message: %s", string(msg_json))
	a.c_imagestream <- []byte(string(msg_json))
}
