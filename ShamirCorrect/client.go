package main

import (
	"encoding/gob"
	"fmt"

	"net"
)

const (
	// Const index of server A (main)
	SERVER_A = 0

	// Const index of server B (secondary)
	SERVER_B = 1

	// Const index of server C (tertiary)
	SERVER_C = 2

	// Const index of server D (???)
	SERVER_D = 3
)

type Client struct {

	// Client data
	P  int
	Id string
	K  int

	// Server connections
	Servers []*net.Conn

	// Encoders
	Encoders []*gob.Encoder

	// Decoders
	Decoders []*gob.Decoder
}

func (client *Client) Init(id string, servers, ports []string, P, K int, bad bool) bool {

	// Grab len
	serverCount := len(servers)

	// If serverCount = 1, copy (Assumption is the IP is the same for all servers)
	if serverCount == 1 {
		servers = []string{servers[0], servers[0], servers[0], servers[0]}
		serverCount = len(servers)
	}

	// If server count is not 4, PANIC (at the disco)
	if serverCount != 4 {
		panic(fmt.Errorf("expected 4 servers IPs but were given %v", serverCount))
	}

	// Verify port count (must be 4 separate)
	if len(ports) != 4 {
		panic(fmt.Errorf("expected 4 servers ports but were given %v", len(ports)))
	}

	// Set identifier
	client.Id = id
	client.P = P
	client.K = K

	// Make arrays
	client.Servers = make([]*net.Conn, 4)
	client.Decoders = make([]*gob.Decoder, 4)
	client.Encoders = make([]*gob.Encoder, 4)

	//Make mate a loop of this later instead.
	// Connect to server A
	connA, encA, decA, typeA, err := ConnectServer(id, servers[0], ports[0])
	if err != nil {
		if !bad {
			panic(err) // Cannot complete protocol when one party is not available
		} else {
			fmt.Printf("[%s] Bad client failed connection (S1) and silently shutting off...", id)
			return false
		}
	}

	// Connect to server B
	connB, encB, decB, typeB, err := ConnectServer(id, servers[1], ports[1])
	if err != nil {
		if !bad {
			panic(err) // Cannot complete protocol when one party is not available
		} else {
			fmt.Printf("[%s] Bad client failed connection (S2) and silently shutting off...", id)
			return false
		}
	}

	// Connect to server C
	connC, encC, decC, typeC, err := ConnectServer(id, servers[2], ports[2])
	if err != nil {
		if !bad {
			panic(err) // Cannot complete protocol when one party is not available
		} else {
			fmt.Printf("[%s] Bad client failed connection (S2) and silently shutting off...", id)
			return false
		}
	}

	// Connect to server D
	connD, encD, decD, typeD, err := ConnectServer(id, servers[3], ports[3])
	if err != nil {
		if !bad {
			panic(err) // Cannot complete protocol when one party is not available
		} else {
			fmt.Printf("[%s] Bad client failed connection (S2) and silently shutting off...", id)
			return false
		}
	}

	// Form arrays
	roles := []int{typeA, typeB, typeC, typeD}
	cons := []*net.Conn{connA, connB, connC, connD}
	encs := []*gob.Encoder{encA, encB, encC, encD}
	decs := []*gob.Decoder{decA, decB, decC, decD}

	// Assign
	allServers := client.AssignServerRole(SERVER_A, roles, cons, encs, decs)
	allServers = allServers && client.AssignServerRole(SERVER_B, roles, cons, encs, decs)
	allServers = allServers && client.AssignServerRole(SERVER_C, roles, cons, encs, decs)
	allServers = allServers && client.AssignServerRole(SERVER_D, roles, cons, encs, decs)

	// Log
	fmt.Printf("[%s] All servers connected: %v\n", client.Id, allServers)

	// Return result of assign
	return allServers

}

func (client *Client) AssignServerRole(role int, roles []int, connections []*net.Conn, encoders []*gob.Encoder, decoders []*gob.Decoder) bool {

	for k, v := range roles {
		if v == role+1 {
			client.Servers[role] = connections[k]
			client.Encoders[role] = encoders[k]
			client.Decoders[role] = decoders[k]
			return true
		}
	}

	// Log failure and return false
	fmt.Printf("[%s] Failed to connect to server %v - Reported roles were: %v.\n", client.Id, role, roles)
	return false

}

func ConnectServer(id, ip, port string) (*net.Conn, *gob.Encoder, *gob.Decoder, int, error) {

	// Connect using TCP, over specified address on specified port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, nil, nil, 0, err
	}

	// Create encoder
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	// Send client join
	e := enc.Encode(Request{RequestType: CLIENTJOIN, Strs: []string{id}})
	if e != nil {
		fmt.Printf("[%s] Error when sending join message: %e", id, e)
	}

	var responseRequest Request
	e = dec.Decode(&responseRequest)
	if e != nil {
		fmt.Printf("[%s] Error when receiving join response: %e", id, e)
	}

	if responseRequest.RequestType != ID {
		fmt.Printf("[%s] Failure when receiving join response - invalid response type.", id)
		return nil, nil, nil, 0, fmt.Errorf("ew")
	}

	// Return base case -> nil, nil
	return &conn, enc, dec, responseRequest.Val1, nil

}

func (client *Client) SendVote(vote int) {

	// Get R1, R2
	r1, r2, r3, r4 := Secrify(vote, client.P, client.K)

	// To array
	shares := []int{r1, r2, r3, r4}

	// Log
	fmt.Printf("[%s] My secret is %v, with R1 = %v, R2 = %v, R3 = %v, and R4 = %v\n", client.Id, vote, r1, r2, r3, r4)

	// Loop over
	for k, v := range shares {

		// Send r1 to S1
		e := client.Encoders[k].Encode(RMessage{Vote: v}.ToRequest())
		if e != nil {
			fmt.Printf("[%s] Error when sending R%v: %e\n", client.Id, k, e)
		}

	}

}

func AwaitResponse(server *net.Conn, dec *gob.Decoder, ch chan Results) {

	// Define
	var res Request

	// Read
	e := dec.Decode(&res)
	if e != nil {
		return
	}

	// Make sure it's a tally
	if res.RequestType != TALLY {
		panic(fmt.Errorf("failed to get tally, found %v request", res.RequestType))
	}

	// Write to channel
	ch <- res.ToTallyMsg()

}

func (client *Client) Shutdown(waitForResults bool) {

	// If wait - we wait for S1, S2, and s3 to return something
	if waitForResults {

		// Create channel
		countChan := make(chan Results, 4)

		// Go wait
		for k := range client.Servers {
			go AwaitResponse(client.Servers[k], client.Decoders[k], countChan)
		}

		// Wait for both to come in (We don't know in which order)
		countA := <-countChan
		countB := <-countChan
		countC := <-countChan
		countD := <-countChan

		// Report if any server detected an error
		if countA.Error || countB.Error || countC.Error || countD.Error {
			fmt.Printf("[%s] One or more servers reported an error while computing tally!\n", client.Id)
		}

		// If agreement, print; otherwise inform of mismatching results.
		if countA.Yes == countB.Yes && countA.No == countB.No && countA.Yes == countC.Yes && countA.No == countC.No && countA.Yes == countD.Yes && countA.No == countD.No {
			fmt.Printf("[%s] Yes Votes: %v, No Votes: %v (Total %v).\n", client.Id, countA.Yes, countA.No, countA.Yes+countA.No)
		} else {
			fmt.Printf("[%s] Received two results that do no agree!\n\tServer A = %+v\n\tServer B = %+v\n\tServer C = %+v\n\tServer D = %+v\n", client.Id, countA, countB, countC, countD)
		}

		// Close channel
		close(countChan)

	}

	// Shutdown
	for _, s := range client.Servers {
		(*s).Close()
	}

}
