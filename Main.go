package main

import (
	"os"
)

var portCounter = 22221

var ip = GetSelfIP()

func main() {

	if len(os.Arguments) > 1 {
		// TODO: Parse args
	} else {
		// Read state from user, manually
	}

}

func createNewClient() {

}

func createNewServer(id string, clientPort string, serverPort string) (server *Server) {
	server = new(Server)
	server.Initialise(id, ip, clientPort, serverPort)
	return
}
