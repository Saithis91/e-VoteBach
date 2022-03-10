package main

const (
	NOTAFUCKINGREQUEST = iota
	JOINNETWORK
	SERVERJOIN
	CLIENTJOIN
	RNUMBER
	ID
)

// Define actual request type
type Request struct {
	RequestType int
	Data        []byte
}

type RMessage struct {
	R_1 int
	R_2 int
	R_3 int
}

type IDMessage struct {
	ID int
}
