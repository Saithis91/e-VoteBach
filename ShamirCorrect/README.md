# Shamir Secret Sharing With Error Detection
This folder contains the implementation for the Basic Shamir Sharing Implementation. So there's no error detection and no error correction.
For this implementation there are ! tests. For a detailed list of instructions for how to run the implementation see [Shamir Basic](../ShamirBasic/README.md).

# Building the Implementation
The implementation can be built into an executable using Go's 'build' command.
```cmd
go build
```

# Running Tests
The tests can be run with the following argument to the executable file.
```cmd
-mode test -i {Test Number}
```

### Test 1
In test 1 a simple 8-voter vote is performed with no corruption from the servers. 
This is a *Deterministic* test.

### Test 2
In test 2 a simple vote where a corrupt server will return a bad R-value. The servers will then Correct the error during the Tally
This is a *Deterministic* test.

### Test 3
In test 3  a simple 8-voter vote is performed with a server, which have 50% chance of being corrupt. If the server becomes corrupt, it will be detected and Corrected during the Tally.
This is a non-*Deterministic* test.

### Test 4
In test 4 a "real world" situation is tested and spawns several clients and decides at random if an error (related to R values) should occur. 
This is a fully non-*Deterministic* test.

### Test 5
In test 5 the correction for $R_i\notin Z_p$ is tested.
This is a *Deterministic* test.
