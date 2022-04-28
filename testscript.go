package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Struct for taking optional arguments to client spanwer
type clientVote struct {
	id            string
	name          string
	ports         []string
	Vote          int
	Seed          int
	P             int
	DoSeed        bool
	SwapPorts     []int
	IgnoreResults bool
	BadMode       int
}

// Self IP address for testing
var localIP = GetSelfIP()

// Slice of spawned proceeses
var db_spawnedProcceses []*os.Process

// Slice of test cases
var testCases = []func() bool{
	RunTest01,
	RunTest02,
	RunTest03,
	RunTest04,
	RunTest05,
	RunTest06,
	RunTest07,
	RunTest08,
	RunTest09,
	RunTest10,
	RunTest11,
	RunTest12,
}

// Dispatches calls
func DispatchTestCall(testID int) {
	if testID-1 < 0 {
		passes := 0
		for k, v := range testCases {
			result := v()
			CheckTest(k+1, result)
			if result {
				passes++
			}
			time.Sleep(time.Millisecond * 500)
		}
		fmt.Println()
		fmt.Printf("--- Test Results: %v/%v passes ---", passes, len(testCases))
		fmt.Println()
		return
	}
	if testID-1 >= len(testCases) {
		fmt.Printf("Unknown test '%v'.\n", testID)
	} else {
		CheckTest(testID, testCases[testID-1]())
	}
}

func CheckTest(id int, res bool) {
	if res {
		fmt.Printf("--- TEST %v PASS ---\n", id)
	} else {
		fmt.Printf("--- TEST %v FAIL ---\n", id)
	}
}

func RunTest01() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 1 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 15, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10001", "-pport", "11001,11002", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10002", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ForthServer", "-port", "10003", "-pport", "11001,11002,11003,11004", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	res := localTestServer.WaitForResults()

	return res.No == 0 && res.Yes == 0
}

func RunTest02() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 2 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10001", []string{"11001"}, []string{localIP, localIP}, 15, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10002", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10003", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "4", "-mode", "server", "-name", "ForthServer", "-port", "10004", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	fmt.Println()
	fmt.Printf("@@@ TEST  2: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{id: "1", name: "yay", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "2", name: "nay", Vote: 0, DoSeed: true, Seed: 1})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 2: Waiting for results\n")
	fmt.Println()
	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 2: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 1 && res.Yes == 1 && !res.Error
}

func RunTest03() bool {

	// Log test
	fmt.Println("--- Running test 3 ---")
	fmt.Println("--- Gauss-Elimination ---")
	fmt.Println()

	// Create equations
	/*coeffs := Matrix{
		{2, 1},
		{-1, 1},
	}

	// Create V-vector
	B := Vector{5, 2}

	// Create matrix
	/*A := AugmentedMatrix(coeffs, B)

	// Solve
	gauss_elim(A)
	X := back_substitute(A)
	fmt.Printf("A:%v\nX:%v\n", A, X)

	return X[0] == 1 && X[1] == 3*/return true

}

func RunTest04() bool {

	// Log test
	fmt.Println("--- Running test 4 ---")
	fmt.Println("--- Lagrange with all variables ---")
	fmt.Println()

	set := []Point{{X: 1, Y: 562}, {X: 2, Y: 1125}, {X: 3, Y: 1687}, {X: 4, Y: 2200}}
	p := Lagrange(0, 1997, set)
	fmt.Printf("P was %v\n", p)

	//Polynomial(set)
	return false
}

func RunTest05() bool {

	// Log test
	fmt.Println("--- Running test 5 ---")
	fmt.Println("--- Gauss-Elimination (P=991) ---")
	fmt.Println()

	// Define field
	/*p := 991

	// Create equations
	coeffs := IntMatrix{
		{2, 1},
		{-1, 1},
	}

	// Create V-vector
	B := IntVector{5, 2}

	// Create matrix
	/*A := AugmentedIntMatrix(coeffs, B)

	// Solve
	gauss_elim_field(A, p)
	X := back_substitute_field(A, p)
	fmt.Printf("A:%v\nX:%v\n", A, X)*/

	return false

}

func RunTest06() bool {

	// Log test
	fmt.Println("--- Running test 6 ---")
	fmt.Println()

	// Create equations simple test
	A := Matrix{
		{2, 1},
		{-1, 1},
	}
	B := Vector{5, 2}
	X := GaussElim(A, B)

	// Test
	AssertIsTrue(X[0] == 1 && X[1] == 3, fmt.Sprintf("[%v, %v] != [1,3]", X[0], X[1]))

	// Create equations simple test
	A = Matrix{
		{6, -8, -1},
		{3, -9, 7},
		{-10, -9, 6},
	}
	B = Vector{97, 156, 56}
	X = GaussElim(A, B)

	// Test
	AssertIsTrue(int(X[0]) == 7 && X[1] == -8 && X[2] == 9, fmt.Sprintf("[%v, %v, %v] != [7,-8,9]", X[0], X[1], X[2]))

	// Create equations simple test
	A = Matrix{
		{1, 1, 1},
		{0, 2, 5},
		{2, 5, -1},
	}
	B = Vector{6, -4, 27}
	X = GaussElim(A, B)

	// Test
	AssertIsTrue(int(X[0]) == 5 && X[1] == 3 && X[2] == -2, fmt.Sprintf("[%v, %v, %v] != [5,3,-2]", X[0], X[1], X[2]))

	// Create equations simple test
	A = Matrix{
		{-6, 2, 8, 0},
		{-3, -1, -8, -7},
		{2, 8, -9, -1},
		{-10, -2, 5, -5},
	}
	B = Vector{-94, 104, 9, -19}
	X = GaussElim(A, B)

	// Test
	AssertIsTrue(int(X[0]) == 5 && X[1] == -8 && X[2] == -6 && X[3] == -9, fmt.Sprintf("[%v, %v, %v, %v] != [5,-8,-6,-9]", X[0], X[1], X[2], X[3]))

	// Create equations simple test
	A = Matrix{
		{0, 0, 0, 1, 1},
		{1, 1, 1, 1, 0},
		{8, 4, 2, 1, 7},
		{27, 9, 3, 1, 13},
		{64, 16, 4, 1, 21},
	}
	B = Vector{1, 0, 14, 39, 84}
	Y := GaussElim(A, B)

	E := Y[len(Y)-1]
	Q := Y

	P := (math.Pow(Q[0], 3)*1 + math.Pow(Q[1], 2)*1 + Q[2]*1 + Q[3]) / (1 - E)

	fmt.Printf("%v\n%v\n", Y, P)

	return true

}

func RunTest07() bool {

	// Log test
	fmt.Println("--- Running test 7 ---")
	fmt.Println("--- Lagrange with all variables ---")
	fmt.Println()

	//Iterate over the points, and see if the CorrectError function, can correct the equation, regardless of which point is corrupt.

	//X=4 is corrupted.
	set := []Point{{X: 1, Y: 563}, {X: 2, Y: 1125}, {X: 3, Y: 1687}, {X: 4, Y: 2200}}

	if punkt, e := CorrectError(set, 1997); e != nil {
		fmt.Printf("\033[31mAn Error Occured, it was: %v\033[37m\n", e)
		return false
	} else {
		fmt.Printf("P was %v\n", punkt)
		AssertIsTrue(punkt == 1, "P was not equal 1 (Case 1)\n")
	}
	//X=3 is corrupted.
	set = []Point{{X: 1, Y: 563}, {X: 2, Y: 1125}, {X: 3, Y: 1686}, {X: 4, Y: 2249}}

	if punkt, e := CorrectError(set, 1997); e != nil {
		fmt.Printf("\033[31mAn Error Occured, it was: %v\033[37m\n", e)
		return false
	} else {
		fmt.Printf("P was %v\n", punkt)
		AssertIsTrue(punkt == 1, "P was not equal 1 (Case 2)\n")
	}
	//X=2 is corrupted.
	set = []Point{{X: 1, Y: 563}, {X: 2, Y: 112}, {X: 3, Y: 1687}, {X: 4, Y: 2249}}

	if punkt, e := CorrectError(set, 1997); e != nil {
		fmt.Printf("\033[31mAn Error Occured, it was: %v\033[37m\n", e)
		return false
	} else {
		fmt.Printf("P was %v\n", punkt)
		AssertIsTrue(punkt == 1, "P was not equal 1 (Case 3)\n")
	}
	//X=1 is corrupted.
	set = []Point{{X: 1, Y: 1}, {X: 2, Y: 1125}, {X: 3, Y: 1687}, {X: 4, Y: 2249}}

	if punkt, e := CorrectError(set, 1997); e != nil {
		fmt.Printf("\033[31mAn Error Occured, it was: %v\033[37m\n", e)
		return false
	} else {
		fmt.Printf("P was %v\n", punkt)
		AssertIsTrue(punkt == 1, "P was not equal 1 (Case 4)\n")
	}

	return true
}

func RunTest08() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 8 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10001", []string{"11001"}, []string{localIP, localIP}, 15, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10002", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10003", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server with dishonest R values.
	if _, e := TestUtil_SpawnTestProcess("-id", "4", "-mode", "server", "-name", "ForthServer-Baddie", "-port", "10004", "-pport", "11001,11002,11003", "-t", "15", "-s", "1", "-b", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	fmt.Println()
	fmt.Printf("@@@ TEST  8: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{id: "1", name: "yay", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "2", name: "nay", Vote: 0, DoSeed: true, Seed: 1})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 8: Waiting for results\n")
	fmt.Println()
	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 8: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 1 && res.Yes == 1 && !res.Error
}

func RunTest09() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 3 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10001", []string{"11001"}, []string{localIP, localIP}, 15, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10002", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10003", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server with dishonest R values.
	if _, e := TestUtil_SpawnTestProcess("-id", "4", "-mode", "server", "-name", "ForthServer", "-port", "10004", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  3: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{id: "1", name: "yay1", Vote: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "2", name: "yay2", Vote: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "3", name: "yay3", Vote: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "4", name: "nay4", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "5", name: "nay5", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "6", name: "nay6", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "7", name: "nay7", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "8", name: "nay8", Vote: 0})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 3: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 3: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 5 && res.Yes == 3 && !res.Error

}

func RunTest10() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 10 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10001", []string{"11001"}, []string{localIP, localIP}, 5, true, 257)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10002", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10003", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server with dishonest R values.
	if _, e := TestUtil_SpawnTestProcess("-id", "4", "-mode", "server", "-name", "ForthServer-Baddie", "-port", "10004", "-pport", "11001,11002,11003", "-t", "15", "-s", "1", "-b", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  3: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{id: "1", name: "yay1", Vote: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "2", name: "yay2", Vote: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "3", name: "yay3", Vote: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "4", name: "nay4", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "5", name: "nay5", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "6", name: "nay6", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "7", name: "nay7", Vote: 0})
	TestUtil_ClientVoteInstance(clientVote{id: "8", name: "nay8", Vote: 0})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 10: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 10: Got results:\n\t%+v\n", res)
	fmt.Printf("@@@ Expected Result Yes: 3 No: 5\n")
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 5 && res.Yes == 3 && !res.Error

}

func RunTest11() bool {

	fmt.Println("--- Running test 11 ---")
	fmt.Println("--- Error Correction on 4 points ---")
	fmt.Println()

	set := []Point{{X: 1, Y: 154}, {X: 2, Y: 48}, {X: 3, Y: 199}, {X: 4, Y: 184}}

	if punkt, e := CorrectError(set, 257); e != nil {
		fmt.Printf("\033[31mAn Error Occured, it was: %v\033[37m\n", e)
		return false
	} else {
		fmt.Printf("P was %v\n", punkt)
		AssertIsTrue(punkt == 1, "P was not equal 1 (Case 1)\n")
	}
	return false
}

func RunTest12() bool {
	// Init rand
	//rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 12 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()
	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10001", []string{"11001"}, []string{localIP, localIP}, 25, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10002", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10003", "-pport", "11001,11002,11003", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server with dishonest R values.
	if _, e := TestUtil_SpawnTestProcess("-id", "4", "-mode", "server", "-name", "ForthServer-Baddie", "-port", "10004", "-pport", "11001,11002,11003", "-t", "15", "-s", "1", "-b", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  10: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Amount to test
	clients := 50 + rand.Intn(587)

	// Spawn voters
	yesVoters := 0
	for i := 0; i < clients; i++ {
		vote := rand.Intn(2)
		name := "nay" + fmt.Sprint(i-yesVoters)
		if vote == 1 {
			yesVoters++
			name = "yay" + fmt.Sprint(yesVoters)
		}
		TestUtil_ClientVoteInstance(clientVote{id: fmt.Sprint(i), name: name, Vote: vote})
	}

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 12: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	// Log how many we expect
	fmt.Printf("We expect %v yes votes and %v no votes.\n", yesVoters, clients-yesVoters)

	fmt.Println()
	fmt.Printf("@@@ TEST 12 Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == clients-yesVoters && res.Yes == yesVoters
}

func AssertIsTrue(condition bool, msg string) {
	if !condition {
		panic(fmt.Errorf("assert condition failed: %s", msg))
	}
}

func TestUtil_ClientVoteInstance(data clientVote) {
	if data.ports == nil {
		data.ports = []string{"10001", "10002", "10003", "10004"}
	}
	if data.SwapPorts == nil {
		data.SwapPorts = []int{0, 1, 2, 3}
	}
	ports := []string{data.ports[data.SwapPorts[0]], data.ports[data.SwapPorts[1]], data.ports[data.SwapPorts[2]], data.ports[data.SwapPorts[3]]}
	args := []string{
		"-mode", "client",
		"-name", data.name,
		"-port", strings.Join(ports, ","),
		"-v", fmt.Sprint(data.Vote),
	}
	if data.DoSeed {
		args = append(args, "-s", fmt.Sprint(data.Seed))
	}
	if !data.IgnoreResults {
		args = append(args, "-w")
	}
	if data.P != 0 {
		args = append(args, "-p", fmt.Sprint(data.P))
	}
	/*if data.BadMode >= 0 {
		args = append(args, "-b", fmt.Sprintf("%v", data.BadMode))
	}*/
	if _, e := TestUtil_SpawnTestProcess(args...); e != nil {
		panic(fmt.Errorf("assert failed: Voter failed to spawn"))
	}
}

func TestUtil_SpawnTestProcess(args ...string) (*exec.Cmd, error) {
	fmt.Printf("[TestUtil] Spawning vote instance with args: %v\n", args)
	proc := exec.Command("./voting", args...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	e := proc.Start()
	if e == nil {
		// Add to process list so we can ensure it is closed
		db_spawnedProcceses = append(db_spawnedProcceses, proc.Process)
	}
	return proc, e
}

func TestUtil_KillTestProcess(args ...string) (*exec.Cmd, error) {
	fmt.Printf("[TestUtil] Killing off Hellspawn: %v\n", args)
	proc := exec.Command("pkill", "voting")
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	e := proc.Start()
	return proc, e
}
