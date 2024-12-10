package network

const (
	LoginMessage    = "LOGIN"
	MovementMessage = "MOVEMENT"
	CombatMessage   = "COMBAT"
	ErrorMessage    = "ERROR"
)

type Message struct {
	Type    string
	Payload string
}
