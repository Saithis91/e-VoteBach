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

	// Create matrix
	A := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// Solve
	gauss_elim(A)
	B := gauss_jordan_elim(A)
	fmt.Printf("A:%v\nB:%v\n", A, B)

	return true

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
