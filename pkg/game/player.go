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
	"sync"
)

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

// map to store players by username
// thread safe! (?)
var Players = struct {
	sync.RWMutex
	data map[string]*Player
}{
	data: make(map[string]*Player),
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
	Players.Lock()
	defer Players.Unlock()
	Players.data[player.Username] = player

	return player
}

func RemovePlayer(username string) {
	Players.Lock()
	defer Players.Unlock()
	delete(Players.data, username)
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
		_, err = fmt.Sscanf(message, "%s %[^\n]", &msg.Type, &msg.Payload)
		if err != nil {
			utils.LogError("Error parsing message:", err)
			continue
		}

		switch msg.Type {
		case messages.MovementMessage:
			var targetX, targetY float64
			n, err := fmt.Sscanf(msg.Payload, "%f %f", &targetX, &targetY)
			if err == nil && n == 2 {
				p.MoveToTarget(targetX, targetY)
			} else {
				fmt.Println("Invalid movement command:", msg.Payload)
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
}

func (p *Player) UpdatePosition(newX, newY float64) {
	if p.X != newX || p.Y != newY {
		p.X = newX
		p.Y = newY
		p.NeedsSaving = true
	}
}
