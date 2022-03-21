package main

// Enum values defining request types
const (
	NOTAFUCKINGREQUEST = iota
	JOINNETWORK
	SERVERJOIN
	CLIENTJOIN
	RNUMBER
	ID
	TALLY
	CLIENTLIST
	INTERSECTION
)

// Define actual request type
type Request struct {
	RequestType int
	Val1        int
	Val2        int
	Strs        []string
}

func (r Request) ToRMsg() RMessage {
	return RMessage{Vote: r.Val1}
}

func (r Request) ToIdMsg() IDMessage {
	return IDMessage{ID: r.Val1}
}

func (r Request) ToTallyMsg() Results {
	return Results{Yes: r.Val1, No: r.Val2}
}

// R-Vote Message (Client -> Server)
type RMessage struct {
	Vote int
}

// Converts the RMessage into a request
func (m RMessage) ToRequest() Request {
	return Request{RequestType: RNUMBER, Val1: m.Vote}
}

// ID Message
type IDMessage struct {
	ID int
}

// Converts the RMessage into a request
func (m IDMessage) ToRequest() Request {
	return Request{RequestType: ID, Val1: m.ID}
}

// Result message (Server -> Client)
type Results struct {
	Yes int
	No  int
}

// Converts the RMessage into a request
func (m Results) ToRequest() Request {
	return Request{RequestType: TALLY, Val1: m.Yes, Val2: m.No}
}
