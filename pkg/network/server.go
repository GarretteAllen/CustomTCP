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
		return fmt.Errorf("could not start the server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server started on", listener.Addr())

	// start game loop in goroutine
	go game.StartGameLoop()

	for {
		conn, err := listener.Accept()
		fmt.Println("Waiting for connections")
		if err != nil {
			utils.LogError("Error accepting connection", err)
			continue
		}
		go s.HandleClient(conn)
	}
}
