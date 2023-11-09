package shared

// Server to client, notifies the start of a round.
type Round struct {
	character rune
}

// Client to server, asks the server to stop the round, or
// sends the user's input if the round had already ended.
type StopRequest struct {
	words map[Category]string
}

// Server to client, notifies all the other clients
// that another user has finished.
type StopNotify struct{}

// Server to client, informs all users the result
// of the last round.
type Result struct {
	scores map[string]int
}

type Category int

const (
	Animal Category = iota
	Profesion
	Pais
	Deporte
	Nombre
	Color
)
