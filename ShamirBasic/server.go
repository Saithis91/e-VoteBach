package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"time"
)

// Struct for a voter instance
type Voter struct {

	// Connection to voter
	Connection *net.Conn

	// ID
	Id string

	// The secret share
	RVal int

	// Gob encoder and deocer
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

//Struct for a partner Server instance
type PartnerServer struct {
	// Connection to Server
	Connection *net.Conn

	// ID
	Id       string
	ServerID uint8

	// Gob encoder and deocer
	Encoder *gob.Encoder
	Decoder *gob.Decoder

	//Common clientList
	commonClientList bool
}

type ConnectionMap map[string]*Voter
type ServerConnectionMap map[string]*PartnerServer

// Struct for server instance
type Server struct {

	// Connection to the partnerConnection
	PartnerConn    *net.Conn
	PartnerEncoder *gob.Encoder

	// Connection to the Pertner Servers
	PartnerConns ServerConnectionMap

	// (Global) map of all Clients connections
	Clientsconnections ConnectionMap

	// Mutex.locks
	mutex *sync.Mutex

	// Name of server (For debugging identification)
	ID       string
	ServerID uint8

	// Self ip
	SelfIP     string   // self IP
	PartnerIPs []string // IP of partner address

	// Port vals
	ListenPort   string   // Port to listen for client/voter input
	PartnerPorts []string // Port to listen and connect to on partners end

	// The time in seconds to vote
	VoteTime int

	// Create channel for tally
	Tally chan Results

	// Self R-value sum
	SelfRSum int

	// Channel for all points (alpha_i, r_i)
	RPoints chan Point

	// The P value
	P int

	// Flag marking if server is main (Handles R1 values)
	MainServer bool

	// Voters shared across servers
	VoterIntersection StringHashSet

	ClientListener *net.Listener
	ServerListener *net.Listener

	// How many servers to expect input from
	serverThresshold int

	// Summed the Votes
	didSum bool
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

	// Set listener
	server.ClientListener = &ln

	// Close connection
	defer (*server.ClientListener).Close()

	// Log we're listening
	fmt.Printf("[%s] Listening on IP and Port: %s\n", server.ID, ln.Addr().String())

	// While running - accept incoming client/voter connections
	for {

		// Accept
		conn, err := ln.Accept()
		if err != nil { // error checking
			return
		}

		// Handle connection
		go server.HandleVoterConnection(&conn)

	}

}

func (server *Server) InitServerSocket(port string) {

	// Determine which server socket to create
	ip := server.SelfIP

	// Begin listening
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}

	// Save listener
	server.ServerListener = &ln

	// Close connection
	defer (*server.ServerListener).Close()

	// Log we're listening
	fmt.Printf("[%s] Listening on IP and Port: %s for other server.\n\n", server.ID, ln.Addr().String())

	// While running - accept incoming partner server connections
	for {

		// Accept
		conn, err := ln.Accept()
		if err != nil { // error checking
			return
		}

		// Handle connection
		PartnerConn := &conn
		PartnerEncoder := gob.NewEncoder(conn)
		go server.HandleServerPartnerConnect(*PartnerConn, *PartnerEncoder)

	}
}

func (server *Server) HandleVoterConnection(conn *net.Conn) {

	decoder := gob.NewDecoder(*conn)

	// Cleans up after connection finish
	defer (*conn).Close()

	// Grab key
	voterAddr := (*conn).RemoteAddr().String()

	// Handle voter/client stuff
	for {
		var newRequest Request
		e := decoder.Decode(&newRequest)
		if e != nil {
			if errors.Is(e, io.EOF) {
				return
			}
		} else {
			switch newRequest.RequestType {
			case CLIENTJOIN:
				server.mutex.Lock()
				voter := Voter{
					Id:         newRequest.Strs[0],
					Connection: conn,
					Encoder:    gob.NewEncoder(*conn),
					Decoder:    decoder,
				}
				server.Clientsconnections[voterAddr] = &voter
				fmt.Printf("[%s] Registered new voter.\n", server.ID)
				voter.Encoder.Encode(Request{RequestType: ID, Val1: int(server.ServerID)})
				server.mutex.Unlock()
				// Would be here where more stuff would be handled like identification, some exchange of keys etc.
			case RNUMBER:
				// As r message
				rm := newRequest.ToRMsg()
				server.mutex.Lock()
				if voter, exists := server.Clientsconnections[voterAddr]; exists {
					voter.RVal = rm.Vote
				} else {
					fmt.Printf("[%s] Unregistered voter attempted to vote!\n", server.ID)
				}
				server.mutex.Unlock()
			}
		}
	}
}

func (server *Server) HandleServerPartnerConnect(conn net.Conn, encoder gob.Encoder) {

	var Pserver PartnerServer
	//Encoder and Decoder
	decoder := gob.NewDecoder(conn)

	//Cleans up after connection finish
	defer (conn).Close()

	// Handle incoming from partner connection
	for {

		var newRequest Request
		e := decoder.Decode(&newRequest)
		if e != nil {
			if errors.Is(e, io.EOF) {
				fmt.Printf("[%s] Connection closed to partner [%s] (EOF).\n", server.ID, Pserver.Id)
				return
			}
		}

		switch newRequest.RequestType {
		case SERVERJOIN:
			server.mutex.Lock()
			sID := newRequest.ToServerJoinMsg().ID
			fmt.Printf("[%s] Connected with partner server with ID: %s.\n", server.ID, sID)
			Pserver = PartnerServer{
				Id:         newRequest.Strs[0],
				ServerID:   uint8(newRequest.Val1),
				Connection: &conn,
				Encoder:    &encoder,
				Decoder:    decoder,
			}
			server.PartnerConns[sID] = &Pserver
			e := encoder.Encode(ServerJoinIDMessage{ID: server.ID, serverID: server.ServerID}.ToResponse())
			if e != nil {
				panic(e)
			}
			server.mutex.Unlock()
			if server.MainServer {
				if server.serverThresshold <= len(server.PartnerConns) {
					go server.waitTime()
				}
			}
		case RNUMBER:
			// We get r-value from partner, and "terminate"
			rm := newRequest.ToRMsg()
			server.mutex.Lock()
			fmt.Printf("[%s] Got a R-tally number from [%s]: %v.\n", server.ID, Pserver.Id, rm.Vote)

			server.RPoints <- Point{X: int(Pserver.ServerID), Y: rm.Vote}
			// Only do EndVotePeriod once
			if !server.didSum {
				server.EndVotePeriod()
				server.didSum = true
			}
			// If we have more RPoints, than the Thresshold, do the final tally.
			if len(server.RPoints) > server.serverThresshold {
				server.DoTally()
			}
			server.mutex.Unlock()
		case CLIENTLIST:
			server.mutex.Lock()
			checklist := CheckmapFromStringSlice(newRequest.Strs)
			common := make([]string, 0)
			// Check Intersection of Clients between 2 servers
			for _, v := range server.Clientsconnections {
				if _, exists := checklist[v.Id]; exists {
					common = append(common, v.Id)
				}
			}
			Pserver.commonClientList = true
			server.VoterIntersection = CheckmapFromStringSlice(common)
			if server.MainServer {
				flag := true
				// If Main server got the Client intersection to all other servers, set flag to True
				for _, p := range server.PartnerConns {
					if !p.commonClientList {
						flag = false
					}
				}
				if flag {
					// if flag = True, start computation of R-values
					if !server.didSum {
						server.EndVotePeriod()
						server.didSum = true
					}
				}
			} else {
				// Send common to other servers.
				if !server.didSum {
					server.sendClients(common)
				}
			}

			server.mutex.Unlock()

		case SERVERRESPONCE:
			// Other server acknowledged the inter connection.
			server.mutex.Lock()
			sID := newRequest.ToServerJoinMsg()
			fmt.Printf("[%s] Got Responce from partner server with ID: %s: %d.\n", server.ID, sID.ID, sID.serverID)
			Pserver = PartnerServer{
				Id:         sID.ID,
				ServerID:   sID.serverID,
				Connection: &conn,
				Encoder:    &encoder,
				Decoder:    decoder,
			}
			server.PartnerConns[sID.ID] = &Pserver
			server.mutex.Unlock()
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

	// Set incoming
	PartnerConn := &conn
	PartnerEncoder := gob.NewEncoder(conn)

	// Send join message
	e := PartnerEncoder.Encode(ServerJoinIDMessage{ID: server.ID, serverID: server.ServerID}.ToRequest())
	if e != nil {
		panic(e)
	}

	// Handle partner connection
	go server.HandleServerPartnerConnect(*PartnerConn, *PartnerEncoder)

	// Return true
	return true

}

func (server *Server) Initialise(serverID int, id, selfIP string, partnerIP []string, listenPort string, partnerPort []string, waitTime int, mainServer bool, prime int) {

	// Init vals
	server.mutex = &sync.Mutex{}
	server.ServerID = uint8(serverID)
	server.ID = id
	server.SelfIP = selfIP
	server.PartnerIPs = partnerIP
	server.Clientsconnections = ConnectionMap{}
	server.PartnerConns = ServerConnectionMap{}
	server.VoteTime = waitTime
	server.Tally = make(chan Results, 1)
	server.RPoints = make(chan Point, 3)
	server.MainServer = mainServer
	server.serverThresshold = 2
	server.didSum = false
	server.P = prime

	// Log what we're doing
	fmt.Printf("[%s][server Startup] I am main: %v\n", id, mainServer)
	fmt.Printf("[%s][Server Startup] Making server for vote-clients at port: %s\n", id, listenPort)
	fmt.Printf("[%s][server Startup] Making connection to %s:%s.\n", id, server.SelfIP, partnerPort)

	// Set port
	server.ListenPort = listenPort
	server.PartnerPorts = partnerPort

	// If serverCount = 1, copy (Assumption is the IP is the same for all servers)
	if len(server.PartnerIPs) == 1 {
		for i := 1; i < len(server.PartnerPorts); i++ {
			server.PartnerIPs = append(server.PartnerIPs, server.PartnerIPs[0])
		}
	}

	//Try connect to partners
	for i := 0; i < len(server.PartnerIPs); i++ {
		if !server.ConnectToServer(server.PartnerIPs[i], server.PartnerPorts[i]) {
			go server.InitServerSocket(server.PartnerPorts[i])
			fmt.Printf("[%s]Could not find other servers", id)
			break
		}

	}

	// Go init server sockets
	go server.InitClientSocket() // socket for clients

}

func (server *Server) WaitForResults() Results {

	// Get results
	results := <-server.Tally
	resultReq := results.ToRequest()

	// Log
	fmt.Printf("[%s] Tally: %v yes vote(s), %v no vote(s), %v total vote(s).\n", server.ID, results.Yes, results.No, results.Yes+results.No)

	// Inform connected clients
	for ip, client := range server.Clientsconnections {
		e := client.Encoder.Encode(resultReq)
		if e != nil {
			fmt.Printf("[%s] Failed to inform client @%s of results.\n", server.ID, ip)
		}
	}

	// terminate
	//(*server.PartnerConn).Close()

	for _, partner := range server.PartnerConns {
		(*partner.Connection).Close()
	}

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

	// End vote period
	//server.EndVotePeriod()
	//Cross reference that clints are the same across servers.
	server.sendClients(server.getClients(server.Clientsconnections))
}

func (server *Server) EndVotePeriod() {

	// Tally up R-values
	server.SelfRSum = 0
	for _, v := range server.Clientsconnections {
		if _, exists := server.VoterIntersection[v.Id]; exists {
			server.SelfRSum = pmod(server.SelfRSum+v.RVal, server.P)
		}
	}

	// Log exit vote period
	fmt.Printf("[%s] Voting period ended. Got R-value of %v\n", server.ID, server.SelfRSum)

	// Put our point into self R-point
	server.RPoints <- Point{X: int(server.ServerID), Y: server.SelfRSum}

	// Send new r-value to partner
	for _, partner := range server.PartnerConns {
		e := partner.Encoder.Encode(RMessage{Vote: server.SelfRSum}.ToRequest())
		if e != nil {
			fmt.Printf("[%s] Failed to send accumulated R-value to partner, %e\n", server.ID, e)
		} else {
			fmt.Printf("[%s] sent accumulated R-value(%d) to partner, %s\n", server.ID, server.SelfRSum, partner.Id)
		}
	}

}

func (server *Server) DoTally() {

	// Grab points
	a := <-server.RPoints
	b := <-server.RPoints
	c := <-server.RPoints

	// Define vars
	var yes_vote, no_vote int

	// Define  array
	points := []Point{a, b, c}
	sort.Sort(PointXSort(points))

	// Log points
	fmt.Printf("[%s] My points for lagrange interpolation is: %v.\n", server.ID, points)

	// Get (yes) votes
	yes_vote = LagrangeXP(0, server.P, points)

	// Get nays
	no_vote = len(server.VoterIntersection) - yes_vote

	// Log in struct
	tally := Results{
		Yes: yes_vote,
		No:  no_vote,
	}

	// Enter into channel
	server.Tally <- tally

}

func (server *Server) Halt() {

	// Close both listeners
	(*server.ClientListener).Close()
	(*server.ServerListener).Close()

}

func (server *Server) getClients(voters ConnectionMap) (strs []string) {
	keys := make([]string, len(voters))
	for _, v := range voters {
		keys = append(keys, v.Id)
	}
	strs = keys
	return
}

func (server *Server) sendClients(input []string) { //RMessage{Vote: server.SelfRSum}.ToRequest()
	for _, partner := range server.PartnerConns {
		e := partner.Encoder.Encode(StringSlice{slice: input}.ToRequest())
		if e != nil {
			fmt.Printf("[%s]  Sending clients %e to %s\n", server.ID, e, partner.Id)
		}
	}
}
