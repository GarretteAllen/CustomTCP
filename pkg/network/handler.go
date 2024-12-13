package network

import (
	"bufio"
	"customtcp/pkg/game"
	"customtcp/pkg/messages"
	"customtcp/pkg/utils"
	"fmt"
	"net"
)

func (s *Server) HandleClient(conn net.Conn) {
	fmt.Println("New connection from:", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	utils.LogInfo("Waiting for username")
	username, err := reader.ReadString('\n')
	if err != nil {
		utils.LogError("Error reading username:", err)
		conn.Close()
		return
	}
	utils.LogInfo("Received username:", username)

	// Remove the newline character from the username
	username = username[:len(username)-1]

	if username == "" {
		utils.LogInfo("Invalid username received")
		errorMsg := messages.Message{
			Type:    messages.ErrorMessage,
			Payload: "Invalid username",
		}
		conn.Write([]byte(fmt.Sprintf("%s %s\n", errorMsg.Type, errorMsg.Payload)))
		conn.Close()
		return
	}

	// Create player after receiving username
	player := game.NewPlayer(conn, username)
	if player == nil {
		fmt.Println("Failed to create player for username:", username)
		conn.Close()
		return
	}

	// welcome message from server to client
	welcomeMsg := messages.Message{
		Type:    messages.WelcomeMessage,
		Payload: fmt.Sprintf("Welcome, %s\n", username),
	}
	fmt.Println("sending to client: ", welcomeMsg.Payload)
	conn.Write([]byte(fmt.Sprintf("%s %s\n", welcomeMsg.Type, welcomeMsg.Payload)))

	player.SendInitialPositions()

	// initial player position
	posMessage := messages.Message{
		Type:    messages.PositionMessage,
		Payload: fmt.Sprintf("%s %.2f, %.2f", username, player.X, player.Y),
	}
	s.BroadcastMessageToAll(fmt.Sprintf("%s %s\n", posMessage.Type, posMessage.Payload))

	// cleanup after disconnect
	defer func() {
		utils.LogInfo("Player '%s' is disconnecting.", username)
		game.RemovePlayer(username)

		disconnectMsg := messages.Message{
			Type:    messages.DisconnectMessage,
			Payload: username,
		}
		s.BroadcastMessageToAll(fmt.Sprintf("%s %s\n", disconnectMsg.Type, disconnectMsg.Payload))

		conn.Close()
	}()

	player.ListenForMessages()
}

func (s *Server) BroadcastMessageToAll(msg string) {
	players := game.PlayersInstance.GetAllPlayers()
	for _, player := range players {
		if player.Conn != nil {
			_, err := player.Conn.Write([]byte(msg))
			if err != nil {
				utils.LogError("Error sending message to player", err)
			}
		}
	}
}
