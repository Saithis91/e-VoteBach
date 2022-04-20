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
	RunTest07,
	RunTest08,
	RunTest09,
	RunTest10,
	RunTest11,
	RunTest12,
	RunTest13,
	RunTest14,
	RunTest15,
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

	// Yields true, does nothing
	return true

	/*
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
		return h0 == 2 */ // Keeping code for documentation

}

func RunTest04() bool {

	// Return true -> does nothing
	return true

	/*
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
	*/ // Keeping code for documentation

}

func RunTest05() bool {

	// Yields true, does nothing
	return true
	/*
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
		return h0_i == h0_s*/ // Keeping code for documentation

}

func RunTest06() bool {

	// Yields true, does nothing
	return true
	/*
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
		return true*/ // Keeping code for documentation!

}

func RunTest07() bool {

	// r1, r2, r3
	r1 := uint8(5)
	r2 := uint8(8)
	r3 := uint8(11)

	// Test D_0(1)
	//LagrangeBasisGf(1, 3, Gf_FromByte(0), []GfPoint{ {X:Gf_FromByte(r1),Y:GF_FromByte(5)} })

	// Calculate L(0)
	l0 := Lagrange0Gf(r1, r2, r3)

	// Print out l0
	fmt.Printf("l0 = %v.\n", l0)

	// Return true
	return l0 == 2

}

func RunTest08() bool {

	// Secrify X=1 for k=1
	r1, r2, r3 := SecrifyGf(1, 1)

	// Log shares
	fmt.Printf("r1: %v r2: %v R3: %v\n", r1, r2, r3)

	// Compute L(0)
	l0 := Lagrange0Gf(r1, r2, r3)

	// Log l0
	fmt.Printf("L(0) = %v\n", l0)

	// L(0) == 1
	return l0 == 1

}

func RunTest09() bool {

	// Secrify X=1 for k=1
	r1, r2, r3 := SecrifyGf(0, 1)

	// Log shares
	fmt.Printf("r1: %v r2: %v R3: %v\n", r1, r2, r3)

	// Compute L(0)
	l0 := Lagrange0Gf(r1, r2, r3)

	// Log l0
	fmt.Printf("L(0) = %v\n", l0)

	// L(0) == 0
	return l0 == 0

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
	fmt.Printf("@@@ TEST  10: Waiting 5s before spawning clients\n")
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
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == 5 && res.Yes == 3
}

func RunTest11() bool {
	// Init rand
	rand.Seed(1)

	// Log test
	fmt.Println("--- Running test 11 ---")
	fmt.Println()

	// Log what we're testing
	fmt.Println("Starting test-server")
	fmt.Println()
	// Create test server
	localTestServer := CreateNewServer(1, "Main Server", "10000", []string{"11001"}, []string{localIP, localIP}, 30, true, 1997)

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
	fmt.Printf("@@@ TEST 11: Waiting for results\n")
	fmt.Println()

	// Wait for local test server
	res := localTestServer.WaitForResults()

	// Log how many we expect
	fmt.Printf("We expect %v yes votes and %v no votes.\n", yesVoters, clients-yesVoters)

	fmt.Println()
	fmt.Printf("@@@ TEST 11 Got results:\n\t%+v\n", res)
	fmt.Println()

	// Wait 1s before passing/failing
	time.Sleep(1 * time.Second)

	// Halt server
	localTestServer.Halt()

	// Do asserts
	return res.No == clients-yesVoters && res.Yes == yesVoters
}

func RunTest12() bool {

	// Secrify X=1 for k=1
	r1, r2, r3 := SecrifyGf(1, 1)

	// Log shares
	fmt.Printf("r1: %v r2: %v R3: %v\n", r1, r2, r3)

	// Compute L(1)=y1, L(2)=y2
	l1 := LagrangeXGf(uint8(1), r1, r2, r3)
	l2 := LagrangeXGf(uint8(2), r1, r2, r3)

	// Compute slope
	m := int(l2-l1) / (2 - 1)

	// Go back one on x-axis
	l0 := int(l1) - m

	// Log l1
	fmt.Printf("L(1) = %v\n", l1)
	// Log l2
	fmt.Printf("L(2) = %v\n", l2)

	// Log l0
	fmt.Printf("L'(0)= %v\n", l0)

	// Log l0
	fmt.Printf("L''(0)= %v\n", int(r1)-int(r2-r1)/(2-1))

	// L(0) == 0
	return l0 == 1

}

func RunTest13() bool {

	// Secrify X=1 for k=1
	r11, r12, r13 := SecrifyGf(1, 1)
	r21, r22, r23 := SecrifyGf(1, 1)

	// Log shares
	fmt.Printf("r11: %v r12: %v R13: %v\n", r11, r12, r13)
	fmt.Printf("r21: %v r22: %v R23: %v\n", r21, r22, r23)

	// As fields...
	gfr11, gfr12, gfr13 := Gf_FromByte(r11), Gf_FromByte(r12), Gf_FromByte(r13)
	gfr21, gfr22, gfr23 := Gf_FromByte(r21), Gf_FromByte(r22), Gf_FromByte(r23)

	// points
	a := GfPoint{X: Gf_FromByte(1), Y: gfr11.Add(gfr21)}
	b := GfPoint{X: Gf_FromByte(2), Y: gfr12.Add(gfr22)}
	c := GfPoint{X: Gf_FromByte(3), Y: gfr13.Add(gfr23)}

	// Compute l0
	l0 := LagrangeGf(uint8(0), a, b, c).ToByte()

	// Compute l0
	fmt.Printf("l0 = %v\n", l0)

	// L(0) == 0
	return l0 == 2

}

func RunTest14() bool {

	// Chosen prime is best prime
	p := 991

	// Secrify X=1 for k=1
	r11, r12, r13 := Secrify(1, p, 1)
	//r11 := 262
	//r12 := 523
	//r13 := 784
	//r21, r22, r23 := Secrify(1, p, 1)

	// Log shares
	fmt.Printf("r11: %v r12: %v R13: %v\n", r11, r12, r13)
	//fmt.Printf("r21: %v r22: %v R23: %v\n", r21, r22, r23)

	// Compute l0
	l0 := LagrangeXP(0, p, []Point{{X: 1, Y: r11}, {X: 2, Y: r12}, {X: 3, Y: r13}})

	// Compute l0
	fmt.Printf("l0 = %v\n", l0)

	// L(0) == 1
	return l0 == 1

}

func RunTest15() bool {

	// Chosen prime is best prime
	p := 991

	// Secrify X=1 for k=1
	r11, r12, r13 := Secrify(1, p, 1)
	r21, r22, r23 := Secrify(1, p, 1)

	// Log shares
	fmt.Printf("r11: %v r12: %v R13: %v\n", r11, r12, r13)
	fmt.Printf("r21: %v r22: %v R23: %v\n", r21, r22, r23)

	// Compute l0
	l0 := LagrangeXP(0, p, []Point{{X: 1, Y: r11 + r21}, {X: 2, Y: r12 + r22}, {X: 3, Y: r13 + r23}})

	// Compute l0
	fmt.Printf("l0 = %v\n", l0)

	// L(0) == 2
	return l0 == 2

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
