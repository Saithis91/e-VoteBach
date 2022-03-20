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

	// Encoders
	EncoderA *gob.Encoder
	EncoderB *gob.Encoder
}

func (client *Client) Init(id, serverIA, serverIB, serverPA, serverPB string, P int) {

	// Set identifier
	client.Id = id
	client.P = P

	// Connect to server A
	connA, encA, err := ConnectServer(id, serverIA, serverPA)
	if err != nil {
		panic(err) // Cannot complete protocol when one party is not available
	}

	// Connect to serverB
	connB, encB, err := ConnectServer(id, serverIB, serverPB)
	if err != nil {
		panic(err) // Cannot complete protocol when one party is not available
	}

	// Set A stuff
	client.ServerA = connA
	client.EncoderA = encA

	// Set B stuff
	client.ServerB = connB
	client.EncoderB = encB

}

func ConnectServer(id, ip, port string) (*net.Conn, *gob.Encoder, error) {

	// Connect using TCP, over specified address on specified port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, nil, err
	}

	// Create encoder
	enc := gob.NewEncoder(conn)

	// Send client join
	e := enc.Encode(Request{RequestType: CLIENTJOIN})
	if e != nil {
		fmt.Printf("[%s] Error when sending join message: %e", id, e)
	}

	// Return base case -> nil, nil
	return &conn, enc, nil

}

func (client *Client) SendVote(vote int) {

	// Get R1, R2
	r1, r2 := Secrify(vote, client.P)

	// Send r1 to S1
	e := client.EncoderA.Encode(RMessage{Vote: r1}.ToRequest())
	if e != nil {
		fmt.Printf("[%s] Error when sending R1: %e", client.Id, e)
	}

	// Send r2 to S2
	e = client.EncoderB.Encode(RMessage{Vote: r2}.ToRequest())
	if e != nil {
		fmt.Printf("[%s] Error when sending R2: %e", client.Id, e)
	}

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
