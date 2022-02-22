package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
)

type ConnectionMap map[string]*net.Conn

type Server struct {

	// (Global) map of all Server connections
	Serverconnections ConnectionMap

	// (Global) map of all Clients connections
	Clientsconnections ConnectionMap

	//Encoder
	//encoder gob.encoder?

	//Decoder
	//decoder gob.decoder?

	// Mutex.locks
	mutex *sync.Mutex

	// Self ip
	IP string

	// Listing port
	Port string

	// counter for ID
	counter int

	//ID for the thread
	myID string
}

func GetSelfIP() string {

	// Dial Google
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	// Make sure the runtime will cleanup after us
	defer conn.Close()

	// Get our connecting address
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Return connecting address str
	return localAddr.IP.String()
}

func (server *Server) InitServerSocket(port string) {

	// Set port
	server.Port = port

	// Begin listening
	ln, err := net.Listen("tcp", server.IP+":"+port)
	if err != nil {
		panic(err)
	}

	// Close somehow
	defer ln.Close()

	// Log we're listening
	fmt.Println("Listening on IP and Port: " + ln.Addr().String())

	// While running - accept incoming connections
	for {

		// Accept
		conn, _ := ln.Accept() // Should do error checking here...
		server.mutex.Lock()
		// Store connection
		server.Serverconnections[conn.RemoteAddr().String()] = &conn
		server.mutex.Unlock()
		// Handle connection
		go server.HandleConnection(&conn)

	}
}

func (server *Server) HandleConnection(conn *net.Conn) {

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

func (server *Server) Initialise(id string, selfIP string, selfPort string) {

	// Init globals
	server.myID = id
	server.mutex = &sync.Mutex{}
	server.Serverconnections = ConnectionMap{}
	server.Clientsconnections = ConnectionMap{}

	// Log what we're doing
	fmt.Println("[Startup Arguments] making server at port:" + selfPort)
	fmt.Printf("[Startup Arguments] Will attempt to connect to %s:%s.\n", server.IP, selfPort)

	server.mutex.Lock()
	// Set IP
	server.myID = selfIP
	// Set port
	server.Port = selfPort
	server.mutex.Unlock()
}
