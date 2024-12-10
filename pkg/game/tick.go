package game

import (
	"time"
)

const tickRate = 30 * time.Millisecond

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
}
