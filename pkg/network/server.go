package network

import (
	"bufio"
	"customtcp/pkg/game"
	"customtcp/pkg/utils"
	"fmt"
	"net"
)

type Server struct {
	address string
	port    int
}

func NewServer(config *utils.Config) *Server {
	return &Server{
		address: config.ServerAddress,
		port:    config.ServerPort,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server started on", listener.Addr())

	// start game loop in goroutine
	go game.StartGameLoop()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	fmt.Println("New connection from:", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	fmt.Println("Waiting for username")
	username, err := reader.ReadString('\n') // Read until newline character
	if err != nil {
		fmt.Println("Error reading username:", err)
		conn.Close()
		return
	}
	fmt.Println("Received username:", username)

	// Remove the newline character from the username
	username = username[:len(username)-1]

	if username == "" {
		fmt.Println("Invalid username received")
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

	positionMessage := fmt.Sprintf("POSITION %.2f %.2f\n", player.X, player.Y)
	_, err = conn.Write([]byte(positionMessage))
	if err != nil {
		fmt.Println("Error sending position to client:", err)
		conn.Close()
		return
	}
	fmt.Println(positionMessage)

	// Listen for further messages from the player
	player.ListenForMessages()

	// cleanup after disconnect
	defer func() {
		fmt.Println("Player", username, "is disconnecting.")
		delete(game.Players, username)
	}()

	defer conn.Close()
}
