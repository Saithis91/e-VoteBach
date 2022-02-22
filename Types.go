package main

const (
	NOTAFUCKINGREQUEST = iota
	JOINNETWORK
	SERVERJOIN
	CLIENTJOIN
)

// Define actual request type
type Request struct {
	RequestType int
	Data        []byte
}
