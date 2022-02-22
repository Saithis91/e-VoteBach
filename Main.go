package main

var portCounter = 22221

var ip = GetSelfIP()

func main() {



}

func createNewClient() {

}

func createNewServer(id string, clientPort string, serverPort string) (server *Server) {
	server = new(Server)
	server.Initialise(id, ip, clientPort, serverPort)
	return
}
