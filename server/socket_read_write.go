package server

import (
	"net"
)

func ReadCommand(conn net.Conn) (string, error) {
	readBuffer := make([]byte, 1024) // read 1KB data
	n, err := conn.Read(readBuffer)
	if err != nil {
		return "", err
	}
	return string(readBuffer[:n]), nil
}

func Respond(conn net.Conn, command string) error {
	_, err := conn.Write([]byte(command))
	if err != nil {
		return err
	}
	return nil
}
