package main

import (
	"fmt"
	"math/big"
	"time"
)

var portCounter = 22221

var ip = GetSelfIP()

func main(){

}


func createNewServer(id string, port string)(server *Server){
	server = new(Server)
	go server.Initialise(id, ip, port)
	return
}