package main

const (
	NOTAFUCKINGREQUEST = iota
	JOINNETWORK
	SERVERJOIN
	CLIENTJOIN
	RNUMBER
)

// Define actual request type
type Request struct {
	RequestType int
	Data        []byte
}

type RMessage struct {
	Vote int
}
