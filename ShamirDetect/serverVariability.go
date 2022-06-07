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

// Sum Behaviours:

// Do honest R-sum
func HonestRSum(server *Server) int {
	// Tally up R-values
	RSum := 0
	for _, v := range server.Clientsconnections {
		if _, exists := server.VoterIntersection[v.Id]; exists {
			RSum = pmod(RSum+v.RVal, server.P)
		}
	}
	return RSum
}

// Do corrupt R-sum based on mode
func CorruptRSumDet(server *Server, mode int) int {
	if mode == 0 {
		fmt.Printf("[BadServer] \033[31mCorrupting sum to P-value: %v.\033[0m\n", server.P)
		return server.P // simply return p
	} else if mode == 1 {
		fmt.Printf("[BadServer] \033[31mCorrupting sum to honest sum +- random offset: %v.\033[0m\n", server.P)
		return HonestRSum(server) + (rand.Intn(server.P*2) - server.P) // Some random offset from honest r-sum (this may be an OK)
	} else if mode == 2 {
		fmt.Printf("[BadServer] \033[31mCorrupting sum to random upper-bounded P-value: %v.\033[0m\n", server.P)
		return rand.Intn(server.P) // random number in field (this may be an OK)
	} else {
		fmt.Printf("[BadServer] \033[31mCorrupting sum to random negative bounded P-value: %v.\033[0m\n", server.P)
		return -rand.Intn(server.P) // Outside of field
	}
}

// Do corrupt R-sum -> pick one of four options
func CorruptRSum(server *Server) int {
	return CorruptRSumDet(server, rand.Intn(4))
}

// Intersection behaviours

// Honest behaviour, performs the intersection
func HonestIntersection(server *Server, input []string) ([]string, bool) {
	checklist := CheckmapFromStringSlice(input)
	common := make([]string, 0)
	if len(checklist) != len(server.Clientsconnections) {
		return common, true
	}
	err := false
	for _, v := range server.Clientsconnections {
		if _, exists := checklist[v.Id]; exists {
			common = append(common, v.Id)
		} else {
			err = true
		}
	}
	return common, err
}

// Corrupt behaviour
func CorruptIntersection(server *Server, input []string) ([]string, bool) {

	mode := rand.Intn(2)
	if mode == 0 {
		fmt.Println("[BadServer] Returned an empty List")
		return make([]string, 0), false
	} else if mode == 1 {
		common, _ := HonestIntersection(server, input)
		size := rand.Intn(len(common))
		common = common[:size]
		fmt.Printf("[BadServer] reduced the ClientList by %v to %v\n", size, common)

		return common, false
	}

	common, _ := HonestIntersection(server, input)
	size := rand.Intn(len(common))
	for i := 0; i < size; i++ {
		common = append(common, fmt.Sprintf("Bogus%v", i))
	}
	fmt.Printf("[BadServer] increased the ClientList by %v to %v\n", size, common)
	return common, false
}
