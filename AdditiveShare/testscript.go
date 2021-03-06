package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// Struct for taking optional arguments to client spanwer
type clientVote struct {
	Id            string
	portA         string
	portB         string
	Vote          int
	Seed          int
	P             int
	DoSeed        bool
	SwapPorts     bool
	IgnoreResults bool
	IsBad         bool
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
		fmt.Printf("\033[32m--- Test Results: %v/%v passes ---\033[0m", passes, len(testCases))
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
		fmt.Printf("\033[32m--- TEST %v PASS ---\033[0m\n", id)
	} else {
		fmt.Printf("\033[31m--- TEST %v FAIL ---\033[0m\n", id)
	}
}

func RunTest01() bool {

	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 1 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server - vote period is 15 seconds.")
	fmt.Println("There will be 2 votes. 1 Yes vote and 1 No vote.")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer("Main Server", "11000", "11001", localIP, 15, true)
	localTestServer.P = 991

	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 1: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{Id: "yay", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{Id: "nay", Vote: 0, DoSeed: true, Seed: 1})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 1: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 1: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 1 && res.Yes == 1

}

func RunTest02() bool {

	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 2 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server - vote period is 15 seconds.")
	fmt.Println("There will be 4 votes. 3 Yes vote and 1 No vote.")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer("Main Server", "11000", "11001", localIP, 15, true)
	localTestServer.P = 991

	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 2: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn yay voter
	TestUtil_ClientVoteInstance(clientVote{Id: "yay1", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{Id: "yay2", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{Id: "yay3", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{Id: "nay1", Vote: 0, DoSeed: true, Seed: 1})

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
	return res.No == 1 && res.Yes == 3

}

func RunTest03() bool {

	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 3 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server - vote period is 15 seconds.")
	fmt.Println("Test 50 clients for test seed = 1 and client seeds are random.")
	fmt.Println("This test will always have 29 yes votes and 21 no votes. The shares will be random.")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer("Main Server", "11000", "11001", localIP, 15, true)
	localTestServer.P = 991

	// Spawn partner server
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 3: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)
	yayvoters := 0
	i := 0
	for i < 50 {
		v := rand.Intn(2)
		yayvoters += v
		// Spawn voter
		TestUtil_ClientVoteInstance(clientVote{Id: fmt.Sprintf("voter %v", (i + 1)), Vote: v})
		i++
	}

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
	return res.No == i-yayvoters && res.Yes == yayvoters

}

func RunTest04() bool {

	// Init Truely rand
	rand.Seed(int64(time.Now().UnixNano()))

	// Log test
	fmt.Println("--- Running test 4 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server - vote period is 15 seconds.")
	fmt.Println("Test 250 voters where everything is random.")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer("Main Server", "11000", "11001", localIP, 30, true)
	localTestServer.P = 991

	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "30", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 4: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)
	yayvoters := 0
	i := 0
	for i < 250 {
		v := rand.Intn(2)
		yayvoters += v
		TestUtil_ClientVoteInstance(clientVote{Id: fmt.Sprintf("voter %v", (i + 1)), Vote: v})
		i++
	}

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 4: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 4: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == i-yayvoters && res.Yes == yayvoters

}

func RunTest05() bool {

	// Init Truely rand
	rand.Seed(int64(time.Now().UnixNano()))

	// Log test
	fmt.Println("--- Running test 5 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server - vote period is 15 seconds.")
	fmt.Println("Test M voters where everything is random.")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer("Main Server", "11000", "11001", localIP, 30, true)
	localTestServer.P = 991

	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "30", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 5: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)
	yayvoters := 0
	yeoldVoters := 251 + rand.Intn(localTestServer.P-251)
	i := 0
	for i < yeoldVoters {
		v := rand.Intn(2)
		yayvoters += v
		TestUtil_ClientVoteInstance(clientVote{Id: fmt.Sprintf("voter %v", (i + 1)), Vote: v})
		i++
	}

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 5: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 5: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == i-yayvoters && res.Yes == yayvoters

}

func RunTest06() bool {

	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 6 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server - vote period is 15 seconds.")
	fmt.Println("There will be 2 voters. 1 Yes voter and 1 No voter who will fail to connect to one server.")
	fmt.Println()

	// Create test server
	localTestServer := CreateNewServer("Main Server", "11000", "11001", localIP, 15, true)
	localTestServer.P = 991

	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 6: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{Id: "yay", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{Id: "nay", Vote: 0, DoSeed: true, Seed: 1, IsBad: true, portB: "46948"})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 6: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 6: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 0 && res.Yes == 1

}

func TestUtil_ClientVoteInstance(data clientVote) {
	var porta, portb string
	if data.portA == "" {
		data.portA = "11000"
	}
	if data.portB == "" {
		data.portB = "11002"
	}
	if !data.SwapPorts {
		porta = data.portA
		portb = data.portB
	} else {
		porta = data.portB
		portb = data.portA
	}
	args := []string{
		"-mode", "client",
		"-id", data.Id,
		"-port.a", porta,
		"-port.b", portb,
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
	if data.IsBad {
		args = append(args, "-b")
	}
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
