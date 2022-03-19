package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
)

type Client struct {

	// Client data
	P  int
	Id string

	// Server connections
	ServerA *net.Conn
	ServerB *net.Conn
}

func (client *Client) Init(id, serverIA, serverIB, serverPA, serverPB string, P int) {

	// Set identifier
	client.Id = id
	client.P = P

	// Connect to server A
	connA, err := ConnectServer(serverIA, serverPA)
	if err != nil {
		panic(err) // Cannot complete protocol when one party is not available
	}

	// Connect to serverB
	connB, err := ConnectServer(serverIB, serverPB)
	if err != nil {
		panic(err) // Cannot complete protocol when one party is not available
	}

	// Set clients
	client.ServerA = connA
	client.ServerB = connB

}

func ConnectServer(ip, port string) (*net.Conn, error) {

	// Connect using TCP, over specified address on specified port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, err
	}

	// Send client join
	gob.NewEncoder(conn).Encode(Request{RequestType: CLIENTJOIN})

	// Return base case -> nil, nil
	return &conn, nil

}

func (client *Client) SendVote(vote int) {

	// Get R1, R2
	r1, r2 := Secrify(vote, client.P)

	// Grab encoders
	s1 := gob.NewEncoder(*client.ServerA)
	s2 := gob.NewEncoder(*client.ServerB)

	// Send r1 to S1
	s1.Encode(RMessage{Vote: r1}.ToRequest())

	// Send r2 to S2
	s2.Encode(RMessage{Vote: r2}.ToRequest())

}

func AwaitResponse(server *net.Conn, ch chan Results) {

	// Create decode
	s := gob.NewDecoder(*server)

	// Define
	res := new(Request)

	// Read
	e := s.Decode(res)
	if e != nil {
		return
	}

	// Make sure it's a tally
	if (*res).RequestType != TALLY {
		panic(fmt.Errorf("failed to get tally, found %v request", (*res).RequestType))
	}

	// Write to channel
	ch <- (*res).ToTallyMsg()

}

func (client *Client) Shutdown(waitForResults bool) {

	// If wait - we wait for S1 or S2 to return something
	if waitForResults {

		// Create channel
		countChan := make(chan Results)

		// Go read
		go AwaitResponse(client.ServerA, countChan)
		go AwaitResponse(client.ServerB, countChan)

		// Wait
		count := <-countChan

		// Log results
		fmt.Printf("Yes Votes: %v, No Votes: %v (Total %v)", count.Yes, count.No, count.Yes+count.No)

		// Close channel
		close(countChan)

	}

	// Shutdown
	(*client.ServerA).Close()
	(*client.ServerB).Close()

}

func Secrify(x, p int) (r1, r2 int) {
	r1 = rand.Intn(p - 1)
	r2 = (x - r1) % p
	return
}
