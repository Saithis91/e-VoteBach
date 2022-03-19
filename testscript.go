package main

import (
	"fmt"
	"os/exec"
	"time"
)

func DispatchTestCall(testID int) {
	if testID <= 0 {
		i := 1
		for i <= 1 {
			DispatchTestCall(i)
		}
	}
	switch testID {
	case 1:
		RunTest01()
	default:
		fmt.Printf("Unknown test '%v'.\n", testID)
	}
}

func CheckTest(id int, res bool) {
	if res {
		fmt.Printf("--- TEST %v PASS ---", id)
	} else {
		fmt.Printf("--- TEST %v FAIL ---", id)
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
	localTestServer := CreateNewServer("Main Server", "11000", "11001", "127.0.0.1", 15)
	localTestServer.P = 991

	// Spawn server
	if _, e := TestUtil_SpawnTestProcess("-mode server -id otherServer -port 11001 -pport 11000 -t 15"); e != nil {
		return false
	}

	// Wait 2s
	time.Sleep(2 * time.Second)

	// Spawn yay voter
	if _, e := TestUtil_SpawnTestProcess("-mode client -id yes -port.a 11000 -port.b 11001 -ip.a 127.0.0.1 -ip.b 127.0.0.1 -v 1"); e != nil {
		return false
	}

	// Spawn nay voter
	if _, e := TestUtil_SpawnTestProcess("-mode client -id yes -port.a 11000 -port.b 11001 -ip.a 127.0.0.1 -ip.b 127.0.0.1 -v 0"); e != nil {
		return false
	}

	// Wait for local test server
	res := localTestServer.WaitForResults()

	// Do asserts
	return res.No == 1 && res.Yes == 1

}

func TestUtil_SpawnTestProcess(args string) (*exec.Cmd, error) {
	proc := exec.Command("./voting", args)
	return proc, proc.Run()
}
