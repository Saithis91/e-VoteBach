package main

import (
	"fmt"
	"math/rand"
)

// Code for malicious server behaviour
const (
	BEHAVIOUR_MODE_CLIENT_INTERSET = 0
	BEHAVIOUR_MODE_WRONG_R_VALUE   = 1
)

func (s *Server) CommonIntersection(fun func(*Server) []string) []string {
	return fun(s) // Variability points
}

func (s *Server) CalculateSum(fun func(*Server) int) int {
	return fun(s) // Variability points
}

// Behaviours:

func HonestRSum(server *Server) int {
	// Tally up R-values
	RSum := 0
	for _, v := range server.Clientsconnections {
		if _, exists := server.VoterIntersection[v.Id]; exists {
			fmt.Printf("[%s] Counting R-vote of %s\n", server.ID, v.Id)
			RSum = RSum + v.RVal
		}
	}
	return RSum
}

func CorruptRSum(server *Server) int {
	mode := rand.Intn(4)
	if mode == 0 {
		fmt.Printf("[BadServer] Corrupting sum to P-value: %v.\n", server.P)
		return server.P // simply return p
	} else if mode == 1 {
		fmt.Printf("[BadServer] Corrupting sum to honest sum +- random offset: %v.\n", server.P)
		return HonestRSum(server) + (rand.Intn(server.P*2) - server.P) // Some random offset from honest r-sum (this may be an OK)
	} else if mode == 2 {
		fmt.Printf("[BadServer] Corrupting sum to random upper-bounded P-value: %v.\n", server.P)
		return rand.Intn(server.P) // random number in field (this may be an OK)
	} else {
		fmt.Printf("[BadServer] Corrupting sum to random negative bounded P-value: %v.\n", server.P)
		return -rand.Intn(server.P) // Outside of field
	}
}

func HonestIntersection(server *Server, input []string) []string {
	checklist := CheckmapFromStringSlice(input)
	common := make([]string, 0)
	for _, v := range server.Clientsconnections {
		if _, exists := checklist[v.Id]; exists {
			common = append(common, v.Id)
		}
	}
	return common
}

func CorruptIntersection(server *Server, input []string) []string {

	mode := rand.Intn(4)
	if mode == 0 {
		fmt.Println("[BadServer] returned Orginal ClientList")
		return input
	} else if mode == 2 {
		fmt.Println("[BadServer] Returned an empty List")
		return make([]string, 0)
	} else if mode == 3 {
		common := HonestIntersection(server, input)
		size := rand.Intn(len(common))
		fmt.Printf("[BadServer] reduced the ClientList by %v\n", size)
		common = common[:size]
		return common
	}

	common := HonestIntersection(server, input)
	size := rand.Intn(len(common))
	fmt.Printf("[BadServer] increased the ClientList by %v\n", size)
	for i := 0; i < size; i++ {
		common = append(common, fmt.Sprint("Bogus%v", i))
	}
	return common

}
