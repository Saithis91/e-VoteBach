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
In test 1 a simple 8-voter vote is performed with no corruption from server or clients. This is a control test to verify detection mechanisms do not give incorrect results on error detection. This is a *Deterministic* test.

### Test 2

In test 2 a simple 8-voter vote is performed with a corrupt server, which will return a incorrect $R-value$, which would result in a miscalculated Vote result.
This is a semi-*Deterministic* test.

### Test 3

In test 3  a simple 8-voter vote is performed with a corrupt server, which will return a incorrect $client list$, which would result in a incorrect Client intersection
This is a semi-*Deterministic* test.

### Test 4

In test 4 a simple 8-voter vote is performed with a server, which will at random either be a Honest or Corrupt server, and the Corrupt server can furthermore differ in which way it will react, either with a $corrupt clientList$, or $R$-value  
This is a non-*Deterministic* test.

### Test 5
In test 5 the outside field case is tested ($R_1\notin Z_p$).
