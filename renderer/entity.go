package renderer

//Entity stuff to Json
type Entity struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Visible bool   `json:"visible"`
}

//EntityCollection stuff to Json
type EntityCollection struct {
	Collection []Entity `json:"collection"`
}
