# Additive Sharing
This folder contains the implementation for additive sharing (2-Server Solution). For this implementation there are 5 tests.

# Building the Implementation
The implementation can be built into an executable using Go's 'build' command.
```cmd
go build
```

# Running Implementation
The implementation compiles to a single executable file. We will list a couple of commands here to get the basics going. The execution can be adjusted on several parameters. A full list of arguments can be found by invoking the exectable with a `-help` flag.
## Servers
The built executable file functions as both the server and client file. To run the serverside (for accepting $r_1$-values), run the executable with arguments:
```cmd
-mode server -m -pip {Partner IP Address} -pport {Partner Port} -port {Listen Port}
```
If the server should be the second server (accpets $r_2$-values) drop the `-m` flag. For localised tests (running on the same machine) the `-pip` argument can also be dropped, as it will then use the local machine's IP.

## Clients
The client takes the arguments:
```cmd
-mode client -id {ClientName} -port.a {Main Server Listen Port} -port.b {Partner Server Listen Port} -v {Voting Value}
```
With the Voting Value $\in\{0,1\}$.

# Running Tests

```cmd
-mode test -i {Test Number}
```
## Linux
The linux distribution also has a makefile that can run the tests with
```cmd
make test t={Test Number}
```

### Test 1
In test 1 a simple 2-voter vote is performed with 1 voter voting yes and 1 voter voting no. The expected result is then a final result of 1 yes vote and 1 no vote.
This is a *Deterministic* test.

### Test 2
In test 2 a simple 4-voter vote is performed with 3 yes votes and 1 no vote. The expected result is then a final result of 3 yes votes and 1 no vote.
This is a *Deterministic* test.

### Test 3
In test 3 a 50-voter vote is performed with 29 yes votes and 21 no votes. This test is using a fixed random generator seed for the servers, but each client is using a unique seed.
This is a semi-*Deterministic* test.

### Test 4
In test 4, a 250-voter vote is performed, but each client is voting randomly, where each vote is $\in\{0,1\}$, and a unique seed.
This is a non-*Deterministic* test. 

### Test 5
In test 5, a $M$-voter vote is performed, for random $M : 251\leq M\leq 990$, but each client is voting randomly, where each vote is $\in\{0,1\}$, and a unique seed.
This is a non-*Deterministic* test. 

### Test 6
In test 6 we perform a 2-voter vote where one voter fails to connect to both servers. We expect the 1 valid client to be counted but the bad client to be dropped from the vote tally.
This is a *Deterministic* test. 

# Windows Powershell
To run a file in Windows Powershell the full path is required (unless the folder is added to the environment variables). So running the first server example on Windows would for example be
```cmd
E:\e-VoteBach\AdditiveShare\voting.exe -mode server -m -pport 11000 -port 11001
```