package network

import (
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
		return fmt.Errorf("Could not start the server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server started on", listener.Addr())

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
	var username string
	_, err := fmt.Fscan(conn, &username)
	if err != nil {
		fmt.Println("Error reading username:", err)
		conn.Close()
		return
	}

	player := game.NewPlayer(conn, username)
	if player == nil {
		fmt.Println("Failed to create player for username:", username)
		conn.Close()
		return
	}

	player.ListenForMessages()
}
