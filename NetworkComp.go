package main

import (
	"log"
	"net"
)

func GetSelfIP() string {

	// Dial Google
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	// Make sure the runtime will cleanup after us
	defer conn.Close()

	// Get our connecting address
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Return connecting address str
	return localAddr.IP.String()
}
