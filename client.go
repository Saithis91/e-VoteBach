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

	// Decoders
	decoderA *gob.Decoder
	decoderB *gob.Decoder
}

func (client *Client) Init(id, serverIA, serverIB, serverPA, serverPB string, P int) {

	// Set identifier
	client.Id = id
	client.P = P

	// Connect to server A
	connA, encA, decA, typeA, err := ConnectServer(id, serverIA, serverPA)
	if err != nil {
		panic(err) // Cannot complete protocol when one party is not available
	}

	// Connect to serverB
	connB, encB, decB, typeB, err := ConnectServer(id, serverIB, serverPB)
	if err != nil {
		panic(err) // Cannot complete protocol when one party is not available
	}

	if typeA == 1 { // Connection A = main

		// Set A stuff
		client.ServerA = connA
		client.EncoderA = encA
		client.decoderA = decA

		// Set B stuff
		client.ServerB = connB
		client.EncoderB = encB
		client.decoderB = decB

	} else if typeB == 1 { // Connection B = main

		// Set A stuff
		client.ServerA = connB
		client.EncoderA = encB
		client.decoderA = decB

		// Set B stuff
		client.ServerB = connA
		client.EncoderB = encA
		client.decoderB = decA

	} else {

		panic(fmt.Errorf("neither server was 'main' server"))

	}

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
	e := enc.Encode(Request{RequestType: CLIENTJOIN})
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
	r1, r2 := Secrify(vote, client.P)

	// Log
	fmt.Printf("[%s] My secret is %v, with R1 = %v and R2 = %v\n", client.Id, vote, r1, r2)

	// Send r1 to S1
	e := client.EncoderA.Encode(RMessage{Vote: r1}.ToRequest())
	if e != nil {
		fmt.Printf("[%s] Error when sending R1: %e\n", client.Id, e)
	}

	// Send r2 to S2
	e = client.EncoderB.Encode(RMessage{Vote: r2}.ToRequest())
	if e != nil {
		fmt.Printf("[%s] Error when sending R2: %e\n", client.Id, e)
	}

}

func AwaitResponse(server *net.Conn, dec *gob.Decoder, ch chan Results) {

	defer recover()

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

	// If wait - we wait for S1 or S2 to return something
	if waitForResults {

		// Log results
		fmt.Printf("[%s] Waiting for results...\n", client.Id)

		// Create channel
		countChan := make(chan Results, 1)

		// Go read
		go AwaitResponse(client.ServerA, client.decoderA, countChan)
		//go AwaitResponse(client.ServerB, client.decoderB, countChan)

		// Wait
		count := <-countChan

		// Log results
		fmt.Printf("[%s] Yes Votes: %v, No Votes: %v (Total %v).\n", client.Id, count.Yes, count.No, count.Yes+count.No)

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
