package game

import (
	"bufio"
	"customtcp/pkg/database"
	"customtcp/pkg/messages"
	"customtcp/pkg/models"
	"customtcp/pkg/utils"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
)

// Define a Players type
type Players struct {
	sync.RWMutex
	data map[string]*Player
}

// Initialize the Players instance
var PlayersInstance = &Players{
	data: make(map[string]*Player),
}

type Player struct {
	Conn        net.Conn
	Username    string
	X, Y        float64
	Attack      int
	Ranged      int
	Hitpoints   int
	Inventory   []models.Item
	NeedsSaving bool
}

func (p *Players) GetAllPlayers() map[string]*Player {
	p.RLock()
	defer p.RUnlock()
	return p.data
}

func NewPlayer(conn net.Conn, username string) *Player {
	playerData, err := database.GetPlayerData(username)
	if err != nil {
		fmt.Println("Error loading player data:", err, "Setting default data")
		playerData = models.Player{
			Username:  username,
			X:         0.0,
			Y:         0.0,
			Attack:    1,
			Ranged:    1,
			Hitpoints: 10,
			Inventory: []models.Item{},
		}

		err := database.SavePlayerData(playerData)
		if err != nil {
			fmt.Println("Error saving new player data:", err)
		} else {
			fmt.Println("New player data saved to database")
		}
	}
	player := &Player{
		Conn:      conn,
		Username:  playerData.Username,
		X:         playerData.X,
		Y:         playerData.Y,
		Attack:    playerData.Attack,
		Ranged:    playerData.Ranged,
		Hitpoints: playerData.Hitpoints,
		Inventory: playerData.Inventory,
	}

	PlayersInstance.Lock()
	defer PlayersInstance.Unlock()
	PlayersInstance.data[player.Username] = player
	fmt.Println("Current players", PlayersInstance.data)
	return player
}

func RemovePlayer(username string) {
	PlayersInstance.Lock()
	defer PlayersInstance.Unlock()
	delete(PlayersInstance.data, username)
	fmt.Printf("Player '%s' removed. Current players: %+v\n", username, PlayersInstance.data)
}

func (p *Player) ListenForMessages() {
	reader := bufio.NewReader(p.Conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from player:", err)
			return
		}
		message = message[:len(message)-1]

		var msg messages.Message
		fmt.Println("received message: ", message)
		parts := strings.Fields(message)
		if len(parts) < 2 {
			utils.LogError("Error parsing message:", err)
			continue
		}

		msg.Type = parts[0]
		msg.Payload = message[len(msg.Type)+1:]

		switch msg.Type {
		case messages.MovementMessage:
			var playerID string
			var targetX, targetY float64
			n, err := fmt.Sscanf(msg.Payload, "%s %f %f", &playerID, &targetX, &targetY)
			if err == nil && n == 3 {
				PlayersInstance.RLock()
				player, exists := PlayersInstance.data[playerID]
				PlayersInstance.RUnlock()

				if exists {
					player.MoveToTarget(targetX, targetY)
					fmt.Println("moving player to", targetX, targetY)
				} else {
					fmt.Println("Invalid movement command")
				}
			}
		case messages.CombatMessage:
			utils.LogInfo("Combat message received: %s", msg.Payload)
		default:
			fmt.Println("Unknown message type:", msg.Type)
		}
	}
}

func (p *Player) MoveToTarget(targetX, targetY float64) {
	deltaX := targetX - p.X
	deltaY := targetY - p.Y

	stepSize := 0.5
	distance := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

	if distance <= stepSize {
		p.X = targetX
		p.Y = targetY
	} else {
		moveRatio := stepSize / distance
		p.X += deltaX * moveRatio
		p.Y += deltaY * moveRatio
	}
	p.UpdatePosition(targetX, targetY)
	fmt.Printf("Player '%s' moved to new position: (%.2f, %.2f)\n", p.Username, p.X, p.Y)

	PlayersInstance.RLock()
	defer PlayersInstance.RUnlock()
	for _, player := range PlayersInstance.data {
		if player.Username != p.Username {
			posMessage := messages.Message{
				Type:    messages.PositionMessage,
				Payload: fmt.Sprintf("%s %.2f %.2f", p.Username, p.X, p.Y),
			}
			_, err := player.Conn.Write([]byte(fmt.Sprintf("%s %s\n", posMessage.Type, posMessage.Payload)))
			if err != nil {
				utils.LogError("Error sending position data to client", err)
			}
		}
	}
	fmt.Printf("Player '%s' moved to new position: (%.2f, %.2f)\n", p.Username, p.X, p.Y)
}

func (p *Player) UpdatePosition(newX, newY float64) {
	if p.X != newX || p.Y != newY {
		p.X = newX
		p.Y = newY
		p.NeedsSaving = true
	}
}

func (p *Player) SendInitialPositions() {
	PlayersInstance.RLock()
	defer PlayersInstance.RUnlock()

	for username, otherPlayer := range PlayersInstance.data {
		if username != p.Username {
			posMessage := messages.Message{
				Type:    messages.PositionMessage,
				Payload: fmt.Sprintf("%s %.2f %.2f", otherPlayer.Username, otherPlayer.X, otherPlayer.Y),
			}
			fmt.Printf("Sending initial position of %s to %s: %s\n", username, p.Username, posMessage.Payload)
			_, err := p.Conn.Write([]byte(fmt.Sprintf("%s %s\n", posMessage.Type, posMessage.Payload)))
			if err != nil {
				fmt.Printf("error sending initial position to %s: %v\n", p.Username, err)
			}
		}
	}
}
