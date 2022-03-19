package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

type ConnectionMap map[string]*net.Conn

type Server struct {

	// Connection to the partnerConnection
	partnerOutgoingConn *net.Conn
	partnerIncomingConn *net.Conn

	// (Global) map of all Clients connections
	Clientsconnections ConnectionMap

	// Mutex.locks
	mutex *sync.Mutex

	// Name of server (For debugging identification)
	ID string

	// Self ip
	SelfIP    string // self IP
	PartnerIP string // IP of partner address

	// Port vals
	ListenPort  string // Port to listen for client/voter input
	PartnerPort string // Port to listen and connect to on partners end

	// The time in seconds to vote
	VoteTime int

	// Create channel for tally
	Tally chan Results

	// R values
	Rs chan int

	// Self R-value sum
	SelfRSum int

	// The P value
	P int
}

func (server *Server) InitServerSocket(asClientSocket bool) {

	// Determine which server socket to create
	var ip, port string
	if asClientSocket {
		ip = server.SelfIP
		port = server.ListenPort
	} else {
		ip = server.SelfIP
		port = server.PartnerPort
	}

	// Begin listening
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}

	// Close connection
	defer ln.Close()

	// Log we're listening
	fmt.Println("Listening on IP and Port: " + ln.Addr().String())

	// While running - accept incoming connections
	for {

		// Accept
		conn, err := ln.Accept()
		if err != nil { // error checking
			panic(err)
		}

		server.mutex.Lock()

		// Handle connection
		if asClientSocket {
			// Store voter connection
			go server.HandleVoterConnection(&conn)
		} else {
			server.partnerIncomingConn = &conn
			go server.HandleServerPartnerConnect()
		}

		server.mutex.Unlock()
	}
}

func (server *Server) HandleVoterConnection(conn *net.Conn) {

	//Decoder
	decoder := gob.NewDecoder(*conn)

	//Cleans up after connection finish
	defer (*conn).Close()

	// Handle voter/client stuff
	for {
		var newRequest Request
		e := decoder.Decode(&newRequest)
		if e != nil {
			fmt.Printf("Error: %e", e)
		} else {

			switch newRequest.RequestType {
			case CLIENTJOIN:
				server.Clientsconnections[(*conn).RemoteAddr().String()] = conn
				fmt.Println("Registered new voter")
				// Would be here where more stuff would be handled like identification, some exchange of keys etc.
			case RNUMBER:
				// As r message
				rm := newRequest.ToRMsg()
				server.Rs <- rm.Vote
			}
		}
	}
}

func (server *Server) HandleServerPartnerConnect() {

	//Encoder and Decoder
	//encoder := gob.NewEncoder(*conn)
	decoder := gob.NewDecoder(*server.partnerIncomingConn)

	//Cleans up after connection finish
	defer (*server.partnerIncomingConn).Close()

	// Do wait period
	go server.waitTime()

	// Handle incoming from partner connection
	for {
		var newRequest Request
		decoder.Decode(&newRequest)

		switch newRequest.RequestType {
		case SERVERJOIN:

		case RNUMBER:
			// We get r-value from partner, and "terminate"
			rm := newRequest.ToRMsg()
			server.DoTally(rm.Vote)
		}
	}

}

func (server *Server) ConnectToServer(ip, port string) {
	// Define address
	addr := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Couldn't reach partner server, Waiting for info.")
		return
	}
	server.partnerIncomingConn = &conn
	go server.HandleServerPartnerConnect()
}

func (server *Server) Initialise(id, selfIP, partnerIP, listenPort, partnerPort string, waitTime int) {

	// Init vals
	server.mutex = &sync.Mutex{}
	server.ID = id
	server.SelfIP = selfIP
	server.PartnerIP = partnerIP
	server.Clientsconnections = ConnectionMap{}

	// Log what we're doing
	fmt.Printf("[Server Startup] Making server for vote-clients at port: %s", listenPort)
	fmt.Printf("[server Startup] Making bi-directional connection to %s:%s.\n", server.SelfIP, partnerPort)

	// Set port
	server.ListenPort = listenPort
	server.PartnerPort = partnerPort

	// Go init server sockets
	go server.InitServerSocket(false) // listen socket
	go server.InitServerSocket(true)  // socket to other server

}

func (server *Server) WaitForResults() Results {

	// Get results
	results := <-server.Tally

	// Log
	fmt.Printf("Tally: %v yes votes, %v no votes, %v total votes.\n", results.Yes, results.No, results.Yes+results.No)

	// TODO: Inform subset of clients

	// terminate
	(*server.partnerIncomingConn).Close()
	(*server.partnerOutgoingConn).Close()

	// Return the results
	return results

}

func (server *Server) waitTime() {

	// Log enter vote period
	fmt.Printf("Server %v has entered voting period.", server.ID)

	// Do wait
	time.Sleep(time.Second * time.Duration(server.VoteTime))

	// Log exit vote period
	fmt.Printf("Server %v has ended voting period. Counting votes...", server.ID)

	// Tally up R-values
	server.SelfRSum = 0
	for v := range server.Rs {
		server.SelfRSum += v
	}

	// Send new r-value to partner
	encoder := gob.NewEncoder(*server.partnerOutgoingConn)
	encoder.Encode(RMessage{Vote: server.SelfRSum})

}

func (server *Server) DoTally(partnerR int) {

	// Get (yes) votes
	yes_vote := (server.SelfRSum + partnerR) % server.P

	// Get nays
	no_vote := len(server.Clientsconnections) - yes_vote

	// Log in struct
	tally := Results{
		Yes: yes_vote,
		No:  no_vote,
	}

	// Enter into channel
	server.Tally <- tally

}
