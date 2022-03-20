package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type ConnectionMap map[string]*net.Conn

type Server struct {

	// Connection to the partnerConnection
	PartnerConn    *net.Conn
	PartnerEncoder *gob.Encoder

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

func (server *Server) InitClientSocket() {

	// Determine which server socket to create
	var ip, port string
	ip = server.SelfIP
	port = server.ListenPort

	// Begin listening
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}

	// Close connection
	defer ln.Close()

	// Log we're listening
	fmt.Printf("[%s] Listening on IP and Port: %s\n", server.ID, ln.Addr().String())

	// While running - accept incoming client/voter connections
	for {

		// Accept
		conn, err := ln.Accept()
		if err != nil { // error checking
			panic(err)
		}

		// Handle connection
		go server.HandleVoterConnection(&conn)

	}
}

func (server *Server) InitServerSocket() {

	// Determine which server socket to create
	var ip, port string
	ip = server.SelfIP
	port = server.PartnerPort

	// Begin listening
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}

	// Close connection
	defer ln.Close()

	// Log we're listening
	fmt.Printf("[%s] Listening on IP and Port: %s for other server.\n\n", server.ID, ln.Addr().String())

	// While running - accept incoming partner server connections
	for {

		// Accept
		conn, err := ln.Accept()
		if err != nil { // error checking
			panic(err)
		}

		// Handle connection
		server.PartnerConn = &conn
		server.PartnerEncoder = gob.NewEncoder(conn)
		go server.HandleServerPartnerConnect()

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
			if errors.Is(e, io.EOF) {
				fmt.Printf("[%s] Connection closed to a voter (EOF).\n", server.ID)
				return
			} else {
				fmt.Printf("[%s] Voter connection error: %e.\n", server.ID, e)
			}
		} else {

			switch newRequest.RequestType {
			case CLIENTJOIN:
				server.mutex.Lock()
				server.Clientsconnections[(*conn).RemoteAddr().String()] = conn
				fmt.Printf("[%s] Registered new voter.\n", server.ID)
				server.mutex.Unlock()
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
	decoder := gob.NewDecoder(*server.PartnerConn)

	//Cleans up after connection finish
	defer (*server.PartnerConn).Close()

	// Handle incoming from partner connection
	for {

		var newRequest Request
		e := decoder.Decode(&newRequest)
		if e != nil {
			if errors.Is(e, io.EOF) {
				fmt.Printf("[%s] Connection closed to partner (EOF).\n", server.ID)
				return
			} else {
				fmt.Printf("[%s] Partner connection error: %e.\n", server.ID, e)
			}
		}

		switch newRequest.RequestType {
		case SERVERJOIN:
			fmt.Printf("[%s] Connected with partner server.\n", server.ID)
			go server.waitTime()
		case RNUMBER:
			// We get r-value from partner, and "terminate"
			rm := newRequest.ToRMsg()
			server.DoTally(rm.Vote)
		}
	}

}

func (server *Server) ConnectToServer(ip, port string) bool {

	// Define address
	target := fmt.Sprintf("%s:%s", ip, port)
	fmt.Printf("[%s] Connecting to : %v \n", server.ID, target)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("[%s] Couldn't reach partner server... Assuming initial server\n", server.ID)
		return false
	}

	// Log other server was reached
	fmt.Printf("[%s] Reached the other server. \n", server.ID)

	// Set incoming
	server.PartnerConn = &conn
	server.PartnerEncoder = gob.NewEncoder(conn)

	// Send join message
	e := server.PartnerEncoder.Encode(Request{RequestType: SERVERJOIN})
	if e != nil {
		panic(e)
	}

	// Handle partner connection
	go server.HandleServerPartnerConnect()

	// Return true
	return true

}

func (server *Server) Initialise(id, selfIP, partnerIP, listenPort, partnerPort string, waitTime int) {

	// Init vals
	server.mutex = &sync.Mutex{}
	server.ID = id
	server.SelfIP = selfIP
	server.PartnerIP = partnerIP
	server.Clientsconnections = ConnectionMap{}
	server.VoteTime = waitTime

	// Log what we're doing
	fmt.Printf("[%s][Server Startup] Making server for vote-clients at port: %s\n", id, listenPort)
	fmt.Printf("[%s][server Startup] Making connection to %s:%s.\n", id, server.SelfIP, partnerPort)

	// Set port
	server.ListenPort = listenPort
	server.PartnerPort = partnerPort

	//Try connect to partner
	if !server.ConnectToServer(server.PartnerIP, server.PartnerPort) {
		go server.InitServerSocket()
	} else {
		go server.waitTime()
	}

	// Go init server sockets
	go server.InitClientSocket() // socket for clients

}

func (server *Server) WaitForResults() Results {

	// Get results
	results := <-server.Tally

	// Log
	fmt.Printf("[%s] Tally: %v yes votes, %v no votes, %v total votes.\n", server.ID, results.Yes, results.No, results.Yes+results.No)

	// TODO: Inform subset of clients

	// terminate
	(*server.PartnerConn).Close()

	// Return the results
	return results

}

func (server *Server) waitTime() {

	// Calculate wait time
	wait := time.Second * time.Duration(server.VoteTime)

	// Log enter vote period
	fmt.Printf("[%s] Entered voting period of %v.\n", server.ID, wait)

	// Do wait
	time.Sleep(wait)

	// Log exit vote period
	fmt.Printf("[%s] Voting period ended. Counting votes...\n", server.ID)

	// Tally up R-values
	server.SelfRSum = 0
	for v := range server.Rs {
		server.SelfRSum += v
	}

	// Send new r-value to partner
	e := server.PartnerEncoder.Encode(RMessage{Vote: server.SelfRSum})
	if e != nil {
		fmt.Printf("[%s] Failed to send accumulated R-value to partner, %e\n", server.ID, e)
	}

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
