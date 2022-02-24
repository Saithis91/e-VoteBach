package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

type ConnectionMap map[string]*net.Conn

type Server struct {

	// (Global) map of all Server connections
	Serverconnections ConnectionMap

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
	counter int

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
		conn, _ := ln.Accept() // Should do error checking here...
		server.mutex.Lock()

		// Handle connection
		if asClientSocket {
			// Store voter connection
			server.Clientsconnections[conn.RemoteAddr().String()] = &conn
			go server.HandleVoterConnection(&conn)
		} else {
			// Store connection
			server.Serverconnections[conn.RemoteAddr().String()] = &conn
			go server.HandleServerPartnerConnect(&conn)
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

func (server *Server) HandleServerPartnerConnect(conn *net.Conn) {

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

func (server *Server) Initialise(id, selfIP, partnerIP, listenPort, partnerPort string) {

	// Init vals
	server.mutex = &sync.Mutex{}
	server.ID = id
	server.SelfIP = selfIP
	server.PartnerIP = partnerIP
	server.Serverconnections = ConnectionMap{}
	server.Clientsconnections = ConnectionMap{}

	// Log what we're doing
	fmt.Printf("[Server Startup] Making server for vote-clients at port: %s", listenPort)
	fmt.Printf("[server Startup] Making bi-directional connection to %s:%s.\n", server.SelfIP, partnerPort)

	server.mutex.Lock()
	// Set port
	server.ListenPort = listenPort
	server.PartnerPort = partnerPort
	server.mutex.Unlock()

}
