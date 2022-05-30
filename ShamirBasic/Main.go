package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var ip = GetSelfIP()

func main() {

	// Inform gob of magic type :D
	gob.Register(Request{})

	var mode, name, partnerPort, partnerIP, portlist, clientIPs string
	var id, testcase, vote, voteperiod, p, k, seed int
	var waitForResults, mainServer, badvariant bool

	flag.StringVar(&mode, "mode", "server", "Specify mode to run with.")
	flag.StringVar(&name, "name", "Turing", "Specify the ID of the instance.")
	flag.StringVar(&partnerIP, "pip", ip, "Specify the IP address of the partner server IP address. Default is localhost.")
	flag.StringVar(&portlist, "port", "11000", "Specify which port to connect to (or listen on if server).")
	flag.StringVar(&partnerPort, "pport", "11001", "Specify which port the connect and listen to as a server.")
	flag.StringVar(&clientIPs, "ip", ip, "Specify which IP to use for servers (seperate with commas, only one address can be specified).")
	flag.IntVar(&id, "id", -1, "Specify the ID of the instance.")
	flag.IntVar(&testcase, "i", -1, "Specify specific test to run. Value <= 0 will run all tests")
	flag.IntVar(&vote, "v", rand.Intn(1-0)+0, "Specify how the client will vote (0/1). Default is false/no (0).")
	flag.IntVar(&voteperiod, "t", 15, "Specify how long the voting period is in seconds.")
	flag.IntVar(&p, "p", 1997, "Specify the prime number to generate secret.")
	flag.IntVar(&k, "k", 1, "Specify the amount of dishonest servers we are preparing for.")
	flag.IntVar(&seed, "s", time.Now().Nanosecond(), "Specify the pseudo-random generator seed.")
	flag.BoolVar(&waitForResults, "w", true, "Specify if client should *NOT* wait for results before terminating server connection.")
	flag.BoolVar(&mainServer, "m", false, "Specify if server Should handle the first part of the secret.")
	flag.BoolVar(&badvariant, "b", false, "Specify if server/client Should behave badly (ignore protocol, crash, etc.).")
	flag.Parse()

	// Init rand
	rand.Seed(int64(seed))

	switch mode {
	case "server":
		server := CreateNewServer(id, name, portlist, strings.Split(partnerPort, ","), strings.Split(partnerIP, ","), voteperiod, mainServer, p)
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
		client := CreateNewClient(name, clientIPs, portlist, p, k, badvariant)
		if client != nil {
			client.SendVote(vote)
			client.Shutdown(waitForResults)
		}
	case "test":
		DispatchTestCall(testcase)
	}

}

func CreateNewClient(id, serverIP, serverPort string, P, K int, bad bool) *Client {

	// Create client
	client := new(Client)
	if client.Init(id, strings.Split(serverIP, ","), strings.Split(serverPort, ","), P, K, bad) {
		// Return client
		return client
	}

	return nil

}

func CreateNewServer(id int, name, listenPort string, parnterPort []string, partnerIP []string, waitTime int, mainServer bool, prime int) *Server {
	server := new(Server)
	server.Initialise(id, name, ip, partnerIP, listenPort, parnterPort, waitTime, mainServer, prime)
	return server
}
