package main

import (
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

	// TODO: Introduction message?

	// Return base case -> nil, nil
	return &conn, nil

}

func (client *Client) SendVote(vote int) {

	// Get R1, R2, R3
	//r1, r2, r3 := Secrify(vote, client.P)

	// Send r1, r3 to S1

	// Send r1, r2 to S2

}

func (client *Client) Shutdown(waitForResults bool) {

}

func Secrify(x, p int) (r1, r2, r3 int) {
	r1 = rand.Intn(p - 1) // TODO: Use actual uniform distribution
	r2 = rand.Intn(p - 1)
	r3 = (x - r1 - r2) % p
	return
}

/*func Verify(r1, r2, r3, x int) bool {
	return r1+r2+r3 == x
}*/
