package main

import (
	"fmt"
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
		fmt.Printf("\033[32m--- TEST %v PASS ---\033[0m\n", id)
	} else {
		fmt.Printf("\033[31m--- TEST %v FAIL ---\033[0m\n", id)
	}
}

func PrintResult(testID int, res Results) {
	fmt.Println()
	fmt.Printf("\033[33m@@@ TEST %v: Got results:\n\t%+v\n\033[0m", testID, res)
	fmt.Println()
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
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10002", "-pport", "11001,11002", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  1: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn voters
	TestUtil_ClientVoteInstance(clientVote{id: "1", name: "yay", Vote: 1, DoSeed: true, Seed: 1})
	TestUtil_ClientVoteInstance(clientVote{id: "2", name: "nay", Vote: 0, DoSeed: true, Seed: 1})

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 1: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	PrintResult(1, res)

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
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10002", "-pport", "11001,11002", "-t", "15", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  2: Waiting 5s before spawning clients\n")
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
	fmt.Printf("@@@ TEST 2: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	PrintResult(2, res)

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 5 && res.Yes == 3
}

func RunTest03() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 03 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()
	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 15, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10001", "-pport", "11001,11002", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10002", "-pport", "11001,11002", "-s", "1"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  03: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Amount to test
	clients := 50

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
	fmt.Printf("@@@ TEST 03: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	// Log how many we expect
	fmt.Printf("We expect %v yes votes and %v no votes.\n", yesVoters, clients-yesVoters)

	PrintResult(03, res)

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == clients-yesVoters && res.Yes == yesVoters
}

func RunTest04() bool {
	// Init Truely rand seed
	rand.Seed(int64(time.Now().UnixNano()))

	// Log test
	fmt.Println("--- Running test 04 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()
	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 20, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10001", "-pport", "11001,11002"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10002", "-pport", "11001,11002"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  04: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Amount to test
	clients := 250

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
	fmt.Printf("@@@ TEST 04: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	// Log how many we expect
	fmt.Printf("We expect %v yes votes and %v no votes.\n", yesVoters, clients-yesVoters)

	PrintResult(04, res)

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == clients-yesVoters && res.Yes == yesVoters
}

func RunTest05() bool {
	// Init Truely rand
	rand.Seed(int64(time.Now().UnixNano()))

	// Log test
	fmt.Println("--- Running test 05 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()
	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 40, true, 1997)

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "2", "-mode", "server", "-name", "otherServer", "-port", "10001", "-pport", "11001,11002"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	time.Sleep(2 * time.Second)
	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-id", "3", "-mode", "server", "-name", "ThirdServer", "-port", "10002", "-pport", "11001,11002"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}
	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST  05: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Amount to test
	lowestThresshold := 1001
	clients := lowestThresshold + rand.Intn(localTestServer.P-lowestThresshold)

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
	fmt.Printf("@@@ TEST 05: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	// Log how many we expect
	fmt.Printf("We expect %v yes votes and %v no votes.\n", yesVoters, clients-yesVoters)

	PrintResult(5, res)

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
		data.ports = []string{"10000", "10001", "10002"}
	}
	if data.SwapPorts == nil {
		data.SwapPorts = []int{0, 1, 2}
	}
	ports := []string{data.ports[data.SwapPorts[0]], data.ports[data.SwapPorts[1]], data.ports[data.SwapPorts[2]]}
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
