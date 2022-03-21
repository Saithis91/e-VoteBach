package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var localIP = GetSelfIP()

func DispatchTestCall(testID int) {
	if testID <= 0 {
		i := 1
		for i <= 2 {
			DispatchTestCall(i)
			i++
		}
	}
	switch testID {
	case 1:
		CheckTest(testID, RunTest01())
	case 2:
		CheckTest(testID, RunTest02())
	case 3:
		CheckTest(testID, RunTest03())
	case 4:
		CheckTest(testID, RunTest04())
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

	// Spawn yay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes", "-port.a", "11000", "-port.b", "11002", "-v", "1", "-s", "1"); e != nil {
		fmt.Print("Yey voter failed\n")
		return false
	}

	// Spawn nay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "no", "-port.a", "11000", "-port.b", "11002", "-v", "0", "-s", "1"); e != nil {
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
	fmt.Printf("@@@ TEST 1: Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

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
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes1", "-port.a", "11000", "-port.b", "11002", "-v", "1", "-s", "1"); e != nil {
		fmt.Print("Yey voter failed\n")
		return false
	}

	// Spawn yay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes2", "-port.a", "11000", "-port.b", "11002", "-v", "1", "-s", "1"); e != nil {
		fmt.Print("Yey voter failed\n")
		return false
	}

	// Spawn yay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes3", "-port.a", "11000", "-port.b", "11002", "-v", "1", "-s", "1"); e != nil {
		fmt.Print("Yey voter failed\n")
		return false
	}

	// Spawn nay voter
	if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "no", "-port.a", "11000", "-port.b", "11002", "-v", "0", "-s", "1"); e != nil {
		fmt.Print("Nay voter failed\n")
		return false
	}

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
	// fmt.Println("There will be 2 votes. 1 Yes vote and 1 No vote.")
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
	fmt.Printf("@@@ TEST 3: Waiting 5s before spawning clients\n")
	fmt.Println()
	time.Sleep(5 * time.Second)
	yayvoters := 0
	i := 0
	for i < 50 {
		v := rand.Intn(2)
		yayvoters += v
		// Spawn yay voter
		if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes", "-port.a", "11000", "-port.b", "11002", "-v", fmt.Sprint(v)); e != nil {
			fmt.Print("A voter failed\n")
			return false
		}
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
	// fmt.Println("There will be 2 votes. 1 Yes vote and 1 No vote.")
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
	for i < 100 {
		v := rand.Intn(2)
		yayvoters += v
		// Spawn yay voter
		if _, e := TestUtil_SpawnTestProcess("-mode", "client", "-id", "yes", "-port.a", "11000", "-port.b", "11002", "-v", fmt.Sprint(v)); e != nil {
			fmt.Print("A voter failed\n")
			return false
		}
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

	// Do asserts
	return res.No == i-yayvoters && res.Yes == yayvoters

}

func TestUtil_SpawnTestProcess(args ...string) (*exec.Cmd, error) {
	fmt.Printf("[TestUtil] Spawning vote instance with args: %v\n", args)
	proc := exec.Command("./voting", args...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	e := proc.Start()
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
