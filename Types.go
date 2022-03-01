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

// R-Vote Message (Client -> Server)
type RMessage struct {
	Vote int
}

// Result message (Server -> Client)
type Results struct {
	Yes int
	No  int
}
