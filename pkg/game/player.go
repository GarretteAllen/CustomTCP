package game

import (
	"bufio"
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

		err := database.SavePlayerData(playerData)
		if err != nil {
			fmt.Println("Error saving new player data:", err)
		} else {
			fmt.Println("New player data saved to database")
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
	reader := bufio.NewReader(p.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from player:", err)
			return
		}
		message = message[:len(message)-1]

		var targetX, targetY float64
		n, err := fmt.Sscanf(message, "MOVE_TO %f %f", &targetX, &targetY)
		if err == nil && n == 2 {
			p.MoveToTarget(targetX, targetY)
			print(targetX, targetY)
		} else {
			fmt.Println("Unknown or invalid movement command:", message)
		}
	}
}

func (p *Player) MoveToTarget(targetX, targetY float64) {
	deltaX := targetX - p.x
	deltaY := targetY - p.y

	stepSize := 0.5
	distance := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

	if distance <= stepSize {
		p.x = targetX
		p.y = targetY
	} else {
		moveRatio := stepSize / distance
		p.x += deltaX * moveRatio
		p.y += deltaY * moveRatio
	}
	fmt.Printf("Player '%s' moved to new position: (%.2f, %.2f)\n", p.username, p.x, p.y)
}
