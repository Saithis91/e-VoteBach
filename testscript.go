package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var localIP = GetSelfIP()

func DispatchTestCall(testID int) {
	if testID <= 0 {
		i := 1
		for i <= 1 {
			DispatchTestCall(i)
		}
	}
	switch testID {
	case 1:
		CheckTest(testID, RunTest01())

	default:
		fmt.Printf("Unknown test '%v'.\n", testID)
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
	if _, e := TestUtil_SpawnTestProcess("-mode", "server", "-id", "otherServer", "-port", "11002", "-pport", "11001", "-t", "15", "-m", "false"); e != nil {
		fmt.Printf("second server failed Error was %v.\n", e)
		return false
	}

	// Wait 2s
	fmt.Println()
	fmt.Printf("@@@ TEST 1: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)

	// Spawn yay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes", "-port.a", "11000", "-port.b", "11002", "-v", "1"); e != nil {
		fmt.Print("Yey voter failed\n")
		return false
	}

	// Spawn nay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "no", "-port.a", "11000", "-port.b", "11002", "-v", "0"); e != nil {
		fmt.Print("Nay voter failed\n")
		return false
	}

	// Wait for results
	fmt.Println()
	fmt.Printf("@@@ TEST 1: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	fmt.Println()
	fmt.Printf("@@@ TEST 1: Got results:\n\t%v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Do asserts
	return res.No == 1 && res.Yes == 1

}

func TestUtil_SpawnTestProcess(args ...string) (*exec.Cmd, error) {
	fmt.Printf("[TestUtil] Spawning vote instance with args: %v\n", args)
	proc := exec.Command("./voting", args...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	e := proc.Start()
	return proc, e
}
