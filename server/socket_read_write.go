package server

import (
	"fmt"
	"net"
)

func ReadCommand(conn net.Conn) (string, error) {
	readBuffer := make([]byte, 1024) // read 1KB data
	//read is a blocking call
	n, err := conn.Read(readBuffer)
	if err != nil {
		return "", err
	}
	if n == 0 {
		fmt.Println("No data read")
		return "", nil
	}
	// Fix: Handle interface{} return from RESPParser
	value, err := RESPParser(readBuffer[:n])
	if err != nil {
		return "", err
	}

	// Type assertion to convert interface{} to string
	if str, ok := value.(string); ok {
		return str, nil
	}
	// Fallback: convert any type to string
	return fmt.Sprintf("%v", value), nil
}

func Respond(conn net.Conn, command string) error {
	_, err := conn.Write([]byte(command))
	if err != nil {
		return err
	}
	return nil
}
