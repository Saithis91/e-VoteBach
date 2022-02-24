package main

import (
	"math/rand"
)

type Client struct {
	P         int
	Id        string
	ServerIPs []string
	Vote      bool
}

func (client *Client) Init(id, serverA, serverB string, vote bool, P int) {

}

func (Client *Client) SendVote() {

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
