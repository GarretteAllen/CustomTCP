package game

import (
	"customtcp/pkg/database"
	"customtcp/pkg/models"
	"fmt"
	"time"
)

const tickRate = 600 * time.Millisecond

func StartGameLoop() {
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			HandleGameTick()
		}
	}
}

func HandleGameTick() {
	// process player movements on tick
	SavePlayerPositions()
}

func SavePlayerPositions() {
	PlayersInstance.RLock()
	defer PlayersInstance.RUnlock()

	for _, player := range PlayersInstance.data {
		if player.NeedsSaving {
			playerData := models.Player{
				Username:  player.Username,
				X:         player.X,
				Y:         player.Y,
				Attack:    player.Attack,
				Ranged:    player.Ranged,
				Hitpoints: player.Hitpoints,
				Inventory: player.Inventory,
			}
			err := database.SavePlayerData(playerData)
			if err != nil {
				fmt.Printf("Failed to save data for player %s: %v\n", player.Username, err)
			} else {
				player.NeedsSaving = false
			}
		}
	}
}
