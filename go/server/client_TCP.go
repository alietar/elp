package server

import (
	"bufio"
	"encoding/json"
	"net"
)

// ReceiveTCP se connecte à un serveur TCP sur l'adresse qu'il veut et décode un JSON dans target
func ReceiveTCP(address string, target any) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Lecture jusqu'au délimiteur (\n)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(line), target)
}

