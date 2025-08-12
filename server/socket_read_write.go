package server

import (
	"fmt"
	"net"
	"redis-internal/core"
	"strings"
)

func ReadCommand(conn net.Conn) (*core.RedisCmd, error) {
	readBuffer := make([]byte, 1024) // read 1KB data
	//read is a blocking call - waits until data arrives
	n, err := conn.Read(readBuffer)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		// fmt.Println("No data read")
		return nil, nil
	}

	fmt.Printf("Raw command received: %q\n", string(readBuffer[:n]))

	// Parse with simplified RESP parser
	tokens, err := core.DecodeCmd(readBuffer[:n])
	if err != nil {
		// fmt.Printf("Decode Cmd failed: %v\n", err)
		return nil, err
	}
	//change the token into redis command format
	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

//

func Respond(conn net.Conn, Command *core.RedisCmd) error {
	// Send RESP Simple String response
	// Send simple OK response
	response := core.EvalAndResponse(Command)
	fmt.Printf("Raw response sent: %q\n", string(response))
	_, err := conn.Write([]byte(response))
	if err != nil {
		return err
	}
	return nil
}
