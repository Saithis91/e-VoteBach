package main

import "fmt"

func DispatchTestCall(testID int) {
	switch testID {
	case 1:
		RunTest01()
	default:
		fmt.Printf("Unknown test '%v'.\n", testID)
	}
}

func RunTest01() {

}
