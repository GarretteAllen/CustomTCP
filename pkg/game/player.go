package game

import (
	"customtcp/pkg/database"
	"customtcp/pkg/models"
	"fmt"
	"math"
	"net"
)

type Player struct {
	conn      net.Conn
	username  string
	x, y      float64
	attack    int
	ranged    int
	hitpoints int
	inventory []models.Item
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
	}
	return &Player{
		conn:      conn,
		username:  playerData.Username,
		x:         playerData.X,
		y:         playerData.Y,
		attack:    playerData.Attack,
		ranged:    playerData.Ranged,
		hitpoints: playerData.Hitpoints,
		inventory: playerData.Inventory,
	}
}

func (p *Player) ListenForMessages() {
	for {
		var message string
		_, err := fmt.Fscan(p.conn, &message)
		if err != nil {
			fmt.Println("Error reading from player:", err)
			return
		}
		var targetX, targetY float64
		n, err := fmt.Sscanf(message, "MOVE_TO %f %f", &targetX, &targetY)
		if err == nil && n == 2 {
			p.MoveToTarget(targetX, targetY)
		} else {
			fmt.Println("Unknown or invalid movement command:", message)
		}
	}
}

func (p *Player) MoveToTarget(targetX, targetY float64) {
	deltaX := targetX - p.x
	deltaY := targetY - p.y

	stepSize := 0.1
	distance := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

	if distance <= stepSize {
		p.x = targetX
		p.y = targetY
	} else {
		moveRatio := stepSize / distance
		p.x += deltaX * moveRatio
		p.y += deltaY * moveRatio
	}
	fmt.Printf("%s moved to position (%.2f, %.2f)\n", p.username, p.x, p.y)
}
