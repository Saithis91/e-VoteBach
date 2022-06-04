# Basic Shamir Secret Sharing
This folder contains the implementation for the Basic Shamir Sharing Implementation. So there's no error detection and no error correction.
For this implementation there are 5 tests. For a detailed list of instructions for how to run the implementation see [Additive Sharing](../AdditiveShare/README.md).

# Building the Implementation
The implementation can be built into an executable using Go's 'build' command.
```cmd
go build
```

# Running The Implementation
The implementation compiles to a single executable file. We will list a couple of commands here to get the basics going. The execution can be adjusted on several parameters (such as the used prime or fix the random seed). A full list of arguments can be found by invoking the exectable with a `-help` flag.
## Servers
The built executable file functions as both the server and client file. To run the serverside, run the executable with arguments:
```cmd
-mode server -id {SID} -pip {S2 IP, S3 IP} -pport {S2 Port, S3 Port} -port {Listen Port}
```
For localised tests (running on the same machine) the `-pip` argument can also be dropped, as it will then use the local machine's IP. The server ID must be a valid ID, $SID\in\{0,1,2\}$. The accepted $r$-value for the server is then $SID+1$.

## Clients
The client takes the arguments:
```cmd
-mode client -id {ClientName} -port "{S1 Listen Port, S2 Listen Port, S3 Listen Port}" -v {Voting Value}
```
With the Voting Value $\in\{0,1\}$. The client will use the local machine's IP when connecting. To specify another IP use `-ip "{S1 ip, S2 ip, S3 ip}"` to specify the IP of a specific server. That is, the server port and ip given with one command are are comma seperated. The order of the servers does not matter as the client will identify the servers on its own.

# Running Tests
The tests can be run with the following argument to the executable file.
```cmd
-mode test -i {Test Number}
```

### Test 1
In test 1 a simple 2-voter vote is performed with 1 voter voting yes and 1 voter voting no. The expected result is then a final result of 1 yes vote and 1 no vote. 
This is a *Deterministic* test.

### Test 2

In test 2 a simple 8-voter vote is performed with 3 yes votes and 5 no vote. The expected result is then a final result of 3 yes votes and 1 no vote. 
This is a *Deterministic* test.

### Test 3

In test 3 a 50-voter vote is performed with 29 yes votes and 21 no votes. This test is using a fixed random generator seed for the servers, but each client is using a unique seed. 
This is a semi-*Deterministic* test.

### Test 4

In test 4, a 250-voter vote is performed, but each client is voting randomly, where each vote is $\in\{0, 1\}, using a unique seed. 
This is a non-*Deterministic* test.

### Test 5

In test 5, a $M$-voter vote is performed, for random $M: 250 \leq M \leq 1996$, but each client is voting randomly, where each vote is  $\in\{0, 1\}, using a unique seed. 
This is a non-*Deterministic* test.
