package server

import (
	"fmt"
	"net"
)

func ReadCommand(conn net.Conn) (string, error) {
	readBuffer := make([]byte, 1024) // read 1KB data
	//read is a blocking call - waits until data arrives
	n, err := conn.Read(readBuffer)
	if err != nil {
		return "", err
	}
	if n == 0 {
		fmt.Println("No data read")
		return "", nil
	}

	fmt.Printf("Raw data received: %q\n", string(readBuffer[:n]))

	// Parse with simplified RESP parser
	value := RESPParser(readBuffer[:n])

	// Convert to string for response
	result := fmt.Sprintf("%v", value)
	fmt.Printf("Final result: %q\n", result)
	return result, nil
}

func Respond(conn net.Conn, command string) error {
	// Send RESP Simple String response
	// Send simple OK response
	response := "+OK\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		return err
	}
	return nil
}
