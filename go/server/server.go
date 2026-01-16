package server

import (
	"encoding/json"
	"log"
	"net"
)

// StartTCPServer démarre un serveur TCP et envoie la structure fournie à chaque client
func StartTCPServer(address string, data any) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Println("TCP server is listening on", address)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Accept error:", err)
				continue
			}
			go handleConnection(conn, data)
		}
	}()

	return nil
}

func handleConnection(conn net.Conn, data any) {
	defer conn.Close()

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	// Délimiteur de message (important côté client)
	payload = append(payload, '\n')

	_, err = conn.Write(payload)
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	log.Printf("Sent %d bytes to %s\n", len(payload), conn.RemoteAddr())
}
