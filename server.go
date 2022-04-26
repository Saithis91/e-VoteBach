package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/rand"
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

//Struct for a partner instance
type PartnerServer struct {
	// Connection to Server
	Connection *net.Conn

	// ID
	Id       string
	ServerID uint8

	// The secret share  // Is this needed?
	//RVal int

	// Gob encoder and deocer
	Encoder *gob.Encoder
	Decoder *gob.Decoder

	//Common clientList
	commonClientList bool

	//Checked ClientList
	comparedClients bool
}

type ConnectionMap map[string]*Voter
type ServerConnectionMap map[string]*PartnerServer

// Function pointers for variability points
type RSumPtr func(*Server) int
type IntersectPtr func(*Server, []string) []string

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

	VoterIntersection StringHashSet

	ClientListener *net.Listener
	ServerListener *net.Listener

	serverThresshold int

	//Summed the Votes
	didSum bool

	//Variable points
	SumCalculation RSumPtr
	IntersectFunc  IntersectPtr
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

	//Cleans up after connection finish
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
			/*if !server.MainServer {
				server.EndVotePeriod()
			}*/
			server.RPoints <- Point{X: int(Pserver.ServerID), Y: rm.Vote}
			if !server.didSum {
				server.EndVotePeriod()
				server.didSum = true
			}
			if len(server.RPoints) >= server.serverThresshold+1 {
				fmt.Printf("[%v] Amount of Points gathered: %v\n Starting Tally\n", server.ID, len(server.RPoints))
				server.DoTally()
			} else {
				fmt.Printf("[%v] Amount of Points gathered: %v\n", server.ID, len(server.RPoints))
			}
			server.mutex.Unlock()
		case CLIENTLIST:
			server.mutex.Lock()
			common := server.IntersectFunc(server, newRequest.Strs)
			Pserver.comparedClients = true
			Pserver.commonClientList = true
			server.VoterIntersection = CheckmapFromStringSlice(common)
			flag := true
			for _, p := range server.PartnerConns {
				if p.comparedClients {
					if !p.commonClientList {
						flag = false
						fmt.Printf("[%v] didn't have common list with %v\n", server.ID, p.Id)
					}
				}
			}
			//Client list across 2 servers wasn't the same
			if !flag {
				//Tell other servers to abort
				server.sendABORT("Non common clientList.")
				// Inform clients of an error occured
				tally := Results{
					Yes:   0,
					No:    0,
					Error: true,
				}
				server.Tally <- tally
				//Client list across all servers was the same.
			} else if server.MainServer {
				// goto next step in process
				if !server.didSum {
					server.EndVotePeriod()
					server.didSum = true
				}
			} else {
				// Send common to other servers
				if !server.didSum {
					server.sendClients(common)
				}
			}

			server.mutex.Unlock()

		case SERVERRESPONCE:
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
		case ABORT:
			server.mutex.Lock()
			sID := newRequest.ToABMsg()

			fmt.Printf("[%v] An ABORT was recieved. Reason %v \n", server.ID, sID)
			// Inform clients of an error occured
			tally := Results{
				Yes:   0,
				No:    0,
				Error: true,
			}
			server.Tally <- tally
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
	server.serverThresshold = 3
	server.Tally = make(chan Results, 1)
	server.RPoints = make(chan Point, server.serverThresshold+1)
	server.MainServer = mainServer
	server.P = prime
	server.SumCalculation = HonestRSum
	server.IntersectFunc = HonestIntersection

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
			fmt.Printf("[%s] Could not find other servers\n", id)
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
	fmt.Printf("[%s] Tally: %v yes vote(s), %v no vote(s), %v total vote(s), Error detected %v.\n", server.ID, results.Yes, results.No, results.Yes+results.No, results.Error)

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

	fmt.Printf("[%s]: clients %s\n", server.ID, server.getClients(server.Clientsconnections))

	// Calculate R sum using specified sum function (Variability point)
	server.SelfRSum = server.SumCalculation(server)

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

func Pop(ints []int, i int) (int, []int) {
	if len(ints) == 1 {
		return ints[0], []int{}
	}
	j := i
	if i == -1 {
		j = rand.Intn(len(ints))
	}
	rval := ints[j] // We must capture return value before returning (Go evaluates multiple returns from right to left...)
	return rval, append(ints[:j], ints[j+1:]...)
}

func argmin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func argMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (server *Server) DoTally() {
	fmt.Printf("[%v] started Tally\n", server.ID)
	// Grab points
	a := <-server.RPoints
	b := <-server.RPoints
	c := <-server.RPoints
	d := <-server.RPoints

	// Define yes vote (it may vary how we get it)
	var yes_vote int

	// Define  array
	points := []Point{a, b, c, d}

	sort.Sort(PointXSort(points))

	// Log points
	fmt.Printf("[%s] My points for lagrange interpolation is: %v.\n", server.ID, points)

	// Pick alpha points given our server ID
	a1, tmp := Pop([]int{0, 1, 2, 3}, int(server.ServerID)-1)
	a2, tmp := Pop(tmp, -1)
	a3, tmp := Pop(tmp, -1)
	a4, _ := Pop(tmp, -1)

	// Define tally object
	var tally Results

	// Define set of sample points
	sample_set := []Point{points[a1], points[a2], points[argmin(a3, a4)]}
	sort.Sort(PointXSort(sample_set))

	// Log points
	fmt.Printf("[%s] My sample points for lagrange interpolation is: %v.\n", server.ID, sample_set)

	// Try compute other point, given selection
	if p := Lagrange(argMax(a3, a4)+1, server.P, sample_set); p != points[argMax(a3, a4)].Y {

		// Log
		fmt.Printf("[%s] Error - Point %v is not a point on polynomium, but %v is. Attempting to correct\n", server.ID, points[a4], p)

		// We now try to recover
		var err error
		if yes_vote, err = CorrectError(points, server.P); err != nil {
			fmt.Printf("[%s] Error - failed to recover from error, %e\n", server.ID, err)

			// Log in struct
			tally = Results{
				Yes:   yes_vote,
				No:    -1,
				Error: true,
			}

			// Enter into channel
			server.Tally <- tally

			return
		}

	} else {

		// Get (yes/yay) votes
		yes_vote = Lagrange(0, server.P, sample_set)

	}

	// Get nays
	no_vote := len(server.VoterIntersection) - yes_vote

	// Log in struct
	tally = Results{
		Yes:   yes_vote,
		No:    no_vote,
		Error: false,
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

func (server *Server) sendClients(input []string) {
	fmt.Printf("[%v] Client list was [%v]\n", server.ServerID, input)
	for _, partner := range server.PartnerConns {
		e := partner.Encoder.Encode(StringSlice{slice: input}.ToRequest())
		if e != nil {
			fmt.Printf("[%s]  Sending clients %e to %s\n", server.ID, e, partner.Id)
		}
	}
}

func (server *Server) sendABORT(reason string) {
	for _, partner := range server.PartnerConns {
		e := partner.Encoder.Encode(ABORTmessage{Message: reason, ServerID: server.ServerID}.ToRequest())
		if e != nil {
			fmt.Printf("[%s]  Sending Abort message %e to %s\n", server.ID, e, partner.Id)
		}
	}
}
