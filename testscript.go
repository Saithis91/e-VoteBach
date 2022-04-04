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
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 15, true)
	localTestServer.P = 991

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

	res := localTestServer.WaitForResults()

	return res.No == 0 && res.Yes == 0
}

func RunTest02() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 6 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()
	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 15, true)
	localTestServer.P = 991

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
	return res.No == 1 && res.Yes == 1
}

func RunTest03() bool {

	// Log purpose
	fmt.Println(" --- Running test 3 --- ")
	fmt.Println(" --- Testing Lagrange --- ")

	// Secret 1 f(0)
	s1 := 3

	// f(x)-coeffs
	fa := []int{2, -1}

	// Secret 2 g(0)
	s2 := -1

	// g(x)-coeffs
	ga := []int{1, 1}

	// Compute shares for f
	r11 := Poly2(1, s1, fa)
	r12 := Poly2(2, s1, fa)
	r13 := Poly2(3, s1, fa)

	// Log f(x)
	fmt.Printf("f(x) = %v, %v, %v.\n", r11, r12, r13)

	// Assert values
	AssertIsTrue(r11 == 4, "R11 != 4")
	AssertIsTrue(r12 == 3, "R12 != 3")
	AssertIsTrue(r13 == 0, "R13 != 0")

	// Compute shares for f
	r21 := Poly2(1, s2, ga)
	r22 := Poly2(2, s2, ga)
	r23 := Poly2(3, s2, ga)

	// Log g(x)
	fmt.Printf("g(x) = %v, %v, %v.\n", r21, r22, r23)

	// Assert values
	AssertIsTrue(r21 == 1, "R21 != 1")
	AssertIsTrue(r22 == 5, "R22 != 5")
	AssertIsTrue(r23 == 11, "R23 != 11")

	// Now do "local sum" (Dont bother asserting these, as they should be correct based on previous)
	h1 := r11 + r21
	h2 := r12 + r22
	h3 := r13 + r23

	// Log h(x)
	fmt.Printf("h(x) = %v, %v, %v.\n", h1, h2, h3)

	// Compute h(x) where x = 0
	h0 := Lagrange(0, Point{X: 1, Y: h1}, Point{X: 2, Y: h2}, Point{X: 3, Y: h3})

	// Log h(0)
	fmt.Printf("h(0) = %v.\n", h0)

	// Return if h2 == 0 (s1 + s2)
	return h0 == 2

}

func RunTest04() bool {

	// Define points
	points := []Point{{X: 1, Y: 5}, {X: 2, Y: 8}, {X: 3, Y: 11}}

	// Compute our implementation
	lag1_impl := LagrangeBasis(0, 0, len(points), points)
	lag2_impl := LagrangeBasis(0, 1, len(points), points)
	lag3_impl := LagrangeBasis(0, 2, len(points), points)

	// Precise compution
	fmt.Printf("WTF IS DIS! = %v\n", ((0 - 2) / (1 - 2) * (0 - 3) / (1 - 3)))

	// Compute slide
	lag1_slid := (0 - 2) / (1 - 2) * (0 - 3) / (1 - 3)
	lag2_slid := (0 - 1) / (2 - 1) * (0 - 3) / (2 - 3)
	lag3_slid := (0 - 1) / (3 - 1) * (0 - 2) / (3 - 2)

	// Log
	fmt.Printf("ell(0) :: Impl = %v; Slide = %v;\n", lag1_impl, lag1_slid)
	fmt.Printf("ell(1) :: Impl = %v; Slide = %v;\n", lag2_impl, lag2_slid)
	fmt.Printf("ell(2) :: Impl = %v; Slide = %v;\n", lag3_impl, lag3_slid)

	// Return compute result
	return lag1_impl == lag1_slid && lag2_impl == lag2_slid && lag3_impl == lag3_slid

}

func RunTest05() bool {

	// Define points
	points := []Point{{X: 1, Y: 5}, {X: 2, Y: 8}, {X: 3, Y: 11}}

	// Define points
	ell0 := LagrangeBasis(0, 0, len(points), points)
	ell1 := LagrangeBasis(0, 1, len(points), points)
	ell2 := LagrangeBasis(0, 2, len(points), points)

	// Log values
	fmt.Printf("ell_0(0) = %v\nell_1(0) = %v\nell_2(0) = %v\n", ell0, ell1, ell2)

	// Compute h(0)
	h0_i := (points[0].Y * ell0) + (points[1].Y * ell1) + (points[2].Y * ell2)
	h0_s := 3*0 + 2 // Slide implementation for h(0)

	// Very basic assert
	fmt.Printf("h_i(0) = %v;\nh_s(0) = %v;\n", h0_i, h0_s)

	// Return compute result
	return h0_i == h0_s

}

func RunTest06() bool {

	// Do 1+1 in Gf(2^8)
	a := Gf_One()
	b := Gf_One()

	// Compute a + b
	c := a.Add(b)

	// Do assert
	AssertIsTrue(c.ToByte() == 0, fmt.Sprintf("%v != 0", c.ToByte()))

	// Do second assert
	c = Gf_FromByte(25).Mul(Gf_FromByte(2))
	AssertIsTrue(c.ToByte() == 50, fmt.Sprintf("%v != 50", c.ToByte()))

	// Do third assert
	c = Gf_FromByte(25).Div(Gf_FromByte(2))
	AssertIsTrue(c.ToByte() == 130, fmt.Sprintf("%v != 130", c.ToByte()))

	// Yey
	return true

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

func TestUtil_KillTestProcess(args ...string) (*exec.Cmd, error) {
	fmt.Printf("[TestUtil] Killing off Hellspawn: %v\n", args)
	proc := exec.Command("pkill", "voting")
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	e := proc.Start()
	return proc, e
}
