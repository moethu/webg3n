package renderer

import "encoding/json"

// Message for client
type Message struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

// SendMessageToClient sends a message to the client
func (app *RenderingApp) SendMessageToClient(action string, value string) {
	m := &Message{Action: action, Value: value}
	msgJSON, err := json.Marshal(m)
	if err != nil {
		app.Application.Log().Error(err.Error())
		return
	}
	app.Log().Info("sending message: " + string(msgJSON))
	app.cImagestream <- []byte(string(msgJSON))
}
