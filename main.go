package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Lire la requête du client
	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// Afficher la requête dans la console du serveur
	fmt.Print("Request received:", string(request))

	// Réponse simple en HTTP
	response := "HTTP/1.1 200 OK\n" +
		"Content-Type: text/plain\n" +
		"Content-Length: 13\n" +
		"\n" + "Hello world!\n"

	// Envoyer la réponse au client
	conn.Write([]byte(response))
}

func main() {
	// Écouter sur le port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on :8080")

	// Boucle pour accepter les connexions
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}

		// Gérer chaque connexion dans une goroutine
		go handleConnection(conn)
	}
}
