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
	Vote int
}

// ID Message
type IDMessage struct {
	ID int
}

// Result message (Server -> Client)
type Results struct {
	Yes int
	No  int
}
