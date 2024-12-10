package main

import (
	"customtcp/pkg/database"
	"customtcp/pkg/network"
	"customtcp/pkg/utils"
	"fmt"
	"log"
)

func main() {
	config, err := utils.LoadConfig("../pkg/configs", "config")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	database.Connect(config.DatabaseURI, config.DatabaseName)

	server := network.NewServer(config)
	if err = server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	fmt.Println("Game server is running")
}
