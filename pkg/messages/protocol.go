package messages

const (
	LoginMessage      = "LOGIN"
	WelcomeMessage    = "WELCOME"
	DisconnectMessage = "DISCONNECT"
	MovementMessage   = "MOVEMENT"
	PositionMessage   = "POSITION"
	CombatMessage     = "COMBAT"
	ErrorMessage      = "ERROR"
)

type Message struct {
	Type    string
	Payload string
}
