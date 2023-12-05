package shared

// Client to server, requests a quick match (no room needed)
type QuickMatch struct {
	name string
}

// Server to client, notifies the start of a round.
type Round struct {
	Character rune
}

// Client to server, asks the server to stop the round, or
// sends the user's input if the round had already ended.
type StopRequest struct {
	Words map[Category]string
}

// Server to client, notifies all the other clients
// that another user has finished.
type StopNotify struct{}

// Server to client, sends all clients the words used during the round,
// organized by their categories.
type Words struct{
	Words map[Category][]string
}

// Client to server, sends the invalid words voted by the client.
type WordsValidation struct{
	Invalid map[Category][]string
}

// Server to client, informs all users the result
// of the last round.
type Score struct {
	Scores map[string]int
}

// Server to client, notifies end of game.
type End struct {}

type Category int

const (
	Animal Category = iota
	Profesion
	Pais
	Deporte
	Nombre
	Color
)
