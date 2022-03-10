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

// R-Vote Message (Client -> Server)
type RMessage struct {
	R_1 int
	R_2 int
	R_3 int
}

type IDMessage struct {
	ID int
}

// Result message (Server -> Client)
type Results struct {
	Yes int
	No  int
}
