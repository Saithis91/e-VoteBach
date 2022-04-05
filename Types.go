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
	SERVERRESPONCE
)

// Define actual request type
type Request struct {
	RequestType int
	Val1        int
	Val2        int
	Val3        int
	Strs        []string
}

func (r Request) ToRMsg() RMessage {
	return RMessage{Vote: uint8(r.Val1)}
}

func (r Request) ToIdMsg() IDMessage {
	return IDMessage{ID: r.Val1}
}

func (r Request) ToTallyMsg() Results {
	return Results{Yes: r.Val1, No: r.Val2}
}

func (r Request) ToStrinceSlice() StringSlice {
	return StringSlice{slice: r.Strs}
}

func (r Request) ToServerJoinMsg() ServerJoinIDMessage {
	return ServerJoinIDMessage{ID: r.Strs[0], serverID: uint8(r.Val1)}
}

// R-Vote Message (Client -> Server)
type RMessage struct {
	Vote uint8
}

// Converts the RMessage into a request
func (m RMessage) ToRequest() Request {
	return Request{RequestType: RNUMBER, Val1: int(m.Vote)}
}

// ID Message
type IDMessage struct {
	ID int
}

// Converts the RMessage into a request
func (m IDMessage) ToRequest() Request {
	return Request{RequestType: ID, Val1: m.ID}
}

// Server Join Message
type ServerJoinIDMessage struct {
	ID       string
	serverID uint8
}

//Converts the ServerJoinIDMessage into a request
func (sID ServerJoinIDMessage) ToRequest() Request {
	return Request{RequestType: SERVERJOIN, Strs: []string{sID.ID}, Val1: int(sID.serverID)}
}

//Converts the ServerJoinIDMessage into a request
func (sID ServerJoinIDMessage) ToResponse() Request {
	return Request{RequestType: SERVERRESPONCE, Strs: []string{sID.ID}, Val1: int(sID.serverID)}
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

type StringSlice struct {
	slice []string
}

// Converts the StringSlice into a request
func (SS StringSlice) ToRequest() Request {
	return Request{RequestType: CLIENTLIST, Strs: SS.slice}
}

// Hash set of strings
type StringHashSet map[string]interface{}

// Convert a string slice into a hash set of strings
func CheckmapFromStringSlice(input []string) StringHashSet {
	result := StringHashSet{}
	for _, v := range input {
		result[v] = nil
	}
	return result
}
