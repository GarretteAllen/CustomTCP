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
	username, err := reader.ReadString('\n') // Read until newline character
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
	conn.Write([]byte(fmt.Sprintf("%s %s\n", welcomeMsg.Type, welcomeMsg.Payload)))

	// initial player position
	posMessage := messages.Message{
		Type:    messages.PositionMessage,
		Payload: fmt.Sprintf("%.2f, %.2f", player.X, player.Y),
	}
	msg := fmt.Sprintf("%s %s\n", posMessage.Type, posMessage.Payload)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		utils.LogError("Error sending position to client", err)
		conn.Close()
		return
	}
	fmt.Println(msg)

	// listen for further messages from the player
	player.ListenForMessages()

	// cleanup after disconnect
	defer func() {
		utils.LogInfo("Player '%s' is disconnecting.", username)
		game.RemovePlayer(username)
	}()

	defer conn.Close()
}
