package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

type Config struct {
	Host string
	Port int
}

func TcpEchoServer(config Config) {
	// Create a socket
	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)

	}
	var concurrent_client int = 0
	/* single threaded TCp client . It will accept only one client at a time and wait till the client closed */
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		concurrent_client++
		// fmt.Printf("Accepet conection: %v concurrent client : %v\n", conn.RemoteAddr(), concurrent_client)
		/* read the command and echo same to the server  continuously till client closed */
		for {
			command, err := ReadCommand(conn)
			if err != nil {
				/* EOF error recieved when client disconnected or close the session */
				if err == io.EOF {
					fmt.Println("clinet Disconnected ", conn.RemoteAddr())
					concurrent_client--
					fmt.Println("Closing the Current connection and ready to accept new client")
					break
				}
				panic(err)
			}
			// fmt.Println("command recived :", command)
			//return the same string to the client
			err = Respond(conn, command)
			if err != nil {
				panic(err)
			}
		}

	}

}
