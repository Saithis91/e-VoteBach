package main

import (
	"flag"
	"fmt"
	"os"
)

var portCounter = 22221

var ip = GetSelfIP()

func main() {

	if len(os.Args) > 1 {

		var mode, id, partnerIP, selfPort, partnerPort string
		var testcase, vote, voteperiod, p int

		flag.StringVar(&mode, "mode", "server", "Specify mode to run with.")
		flag.StringVar(&id, "id", "Turing", "Specify the ID of the instance.")
		flag.StringVar(&partnerIP, "pip", ip, "Specify the IP address of the partner server IP address. Default is localhost.")
		flag.StringVar(&selfPort, "port", "11", "Specify which port to connect to as client or listen to if server.")
		flag.StringVar(&partnerPort, "pport", "11001", "Specify which port the connect and listen to as a server.")
		flag.IntVar(&testcase, "i", -1, "Specify specific test to run.")
		flag.IntVar(&vote, "v", 0, "Specify how the client will vote (0/1). Default is false/no (0).")
		flag.IntVar(&voteperiod, "t", 15, "Specify how long the voting period is in seconds.")
		flag.IntVar(&p, "p", 991, "Specify the prime number to generate secret (Default is 991).")
		flag.Parse()

		switch mode {
		case "server":
			server := createNewServer(id, selfPort, partnerPort, partnerIP)
			server.VoteTime = voteperiod
			break
		case "client":
			if vote < 0 || vote > 1 {
				fmt.Println("Invalid vote. Must be an integer value of 0 or 1.")
				return
			}
			if p < 0 {
				fmt.Println("Invalid P-value. Must be greater than 0 (and prime).")
			}
			break
		case "test":
			DispatchTestCall(testcase)
			break
		case "help":
			flag.PrintDefaults()
			return
		}

	} else {
		// Read state from user, manually
	}

}

func createNewClient(id, serverIP, serverPort string) {

}

func createNewServer(id, listenPort, parnterPort, partnerIP string) *Server {
	server := new(Server)
	server.Initialise(id, ip, partnerIP, listenPort, parnterPort)
	return server
}
