package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
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

	// counter for ID
	//counter int

	// The time in seconds to vote
	VoteTime int
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
			server.Clientsconnections[conn.RemoteAddr().String()] = &conn
			go server.HandleVoterConnection(&conn)
		} else {
			server.partnerIncomingConn = &conn
			go server.HandleServerPartnerConnect()
		}

		server.mutex.Unlock()
	}
}

func (server *Server) HandleVoterConnection(conn *net.Conn) {

	//Encoder and Decoder
	//encoder := gob.NewEncoder(*conn)
	decoder := gob.NewDecoder(*conn)

	//Cleans up after connection finish
	defer (*conn).Close()

	for {
		var newRequest Request
		decoder.Decode(&newRequest)

		switch newRequest.RequestType {

		}
	}
}

func (server *Server) HandleServerPartnerConnect() {

	//Encoder and Decoder
	//encoder := gob.NewEncoder(*conn)
	decoder := gob.NewDecoder(*server.partnerIncomingConn)

	//Cleans up after connection finish
	defer (*server.partnerIncomingConn).Close()

	for {
		var newRequest Request
		decoder.Decode(&newRequest)

		switch newRequest.RequestType {
		case SERVERJOIN:

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
	server.partnerOutgoingConn = &conn
}

func (server *Server) Initialise(id, selfIP, partnerIP, listenPort, partnerPort string) {

	// Init vals
	server.mutex = &sync.Mutex{}
	server.ID = id
	server.SelfIP = selfIP
	server.PartnerIP = partnerIP
	server.Clientsconnections = ConnectionMap{}

	// Log what we're doing
	fmt.Printf("[Server Startup] Making server for vote-clients at port: %s", listenPort)
	fmt.Printf("[server Startup] Making bi-directional connection to %s:%s.\n", server.SelfIP, partnerPort)

	server.mutex.Lock()
	// Set port
	server.ListenPort = listenPort
	server.PartnerPort = partnerPort
	go server.InitServerSocket(false)
	go server.InitServerSocket(true)
	server.mutex.Unlock()

}
