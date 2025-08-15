package main

import (
	"flag"
	"fmt"

	"redis-internal/server"
)

var config struct {
	Host string
	Port int
}

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the dice server")
	flag.IntVar(&config.Port, "port", 7379, "port for the dice server")
	flag.Parse()
}

func main() {
	setupFlags()
	fmt.Println("Starting the NiniDB server...")

	// Convert to server.Config type and call the function
	serverConfig := server.Config{
		Host: config.Host,
		Port: config.Port,
	}
	//server.TcpEchoServer(serverConfig)
	server.RunAsyncTCPServer(serverConfig)
}
