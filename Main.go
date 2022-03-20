package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var ip = GetSelfIP()

func main() {

	// Inform gob of magic type :D
	gob.Register(Request{})

	var mode, id, partnerIP, selfPort, partnerPort, aPort, bPort, aIp, bIp string
	var testcase, vote, voteperiod, p, seed int
	var waitForResults, mainServer bool

	flag.StringVar(&mode, "mode", "server", "Specify mode to run with.")
	flag.StringVar(&id, "id", "Turing", "Specify the ID of the instance.")
	flag.StringVar(&partnerIP, "pip", ip, "Specify the IP address of the partner server IP address. Default is localhost.")
	flag.StringVar(&selfPort, "port", "11000", "Specify which port to connect to as client or listen to if server.")
	flag.StringVar(&partnerPort, "pport", "11001", "Specify which port the connect and listen to as a server.")
	flag.StringVar(&aPort, "port.a", "11000", "Specify which port to use for server A.")
	flag.StringVar(&bPort, "port.b", "11001", "Specify which port to use for server B.")
	flag.StringVar(&aIp, "ip.a", ip, "Specify which IP to use for server A.")
	flag.StringVar(&bIp, "ip.b", ip, "Specify which IP to use for server B.")
	flag.IntVar(&testcase, "i", -1, "Specify specific test to run. Value <= 0 will run all tests")
	flag.IntVar(&vote, "v", 0, "Specify how the client will vote (0/1). Default is false/no (0).")
	flag.IntVar(&voteperiod, "t", 15, "Specify how long the voting period is in seconds.")
	flag.IntVar(&p, "p", 991, "Specify the prime number to generate secret.")
	flag.IntVar(&seed, "s", time.Now().Nanosecond(), "Specify the pseudo-random generator seed.")
	flag.BoolVar(&waitForResults, "w", true, "Specify if client should wait for results before terminating server connection.")
	flag.BoolVar(&mainServer, "m", false, "Specify if server Should handle the first part of the secret.")
	flag.Parse()

	// Init rand
	rand.Seed(int64(seed))

	switch mode {
	case "server":
		server := CreateNewServer(id, selfPort, partnerPort, partnerIP, voteperiod, mainServer)
		server.P = p
		server.WaitForResults()
	case "client":
		if vote < 0 || vote > 1 {
			fmt.Println("Invalid vote. Must be an integer value of 0 or 1.")
			return
		}
		if p <= 3 { // A protocol for secure addition, page 13
			fmt.Println("Invalid P-value. Must be greater than 3 (and prime).")
			return
		}
		client := CreateNewClient(id, aIp, aPort, bIp, bPort, p)
		client.SendVote(vote)
		client.Shutdown(waitForResults)
	case "test":
		DispatchTestCall(testcase)

	}

}

func CreateNewClient(id, serverIPA, serverPortA, serverIPB, serverPortB string, P int) *Client {

	// Create client
	client := new(Client)
	client.Init(id, serverIPA, serverIPB, serverPortA, serverPortB, P)

	// Return client
	return client

}

func CreateNewServer(id, listenPort, parnterPort, partnerIP string, waitTime int, mainServer bool) *Server {
	server := new(Server)
	server.Initialise(id, ip, partnerIP, listenPort, parnterPort, waitTime, mainServer)
	return server
}
