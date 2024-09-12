package main

import (
	"fmt"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("Callers from main:")
	if err := DumpCallerStack(); err != nil {
		fmt.Printf("Error dumping caller stack: %v\n", err)
	}
	func2()
}

func func2() {
	fmt.Println("\nCallers from func2:")
	if err := DumpCallerStack(); err != nil {
		fmt.Printf("Error dumping caller stack: %v\n", err)
	}
}

func DumpCallerStack() error {
	const maxDepth = 100
	callers := make([]string, 0, maxDepth) // Pre-allocate slice for efficiency

	for i := 1; i < maxDepth; i++ {
		caller, err := getCallerName(i)
		if err != nil {
			if err == ErrCallerNotFound {
				break // End of call stack reached
			}
			return fmt.Errorf("error getting caller name: %w", err)
		}
		callers = append(callers, caller)
	}

	// Print callers in reverse order (most recent call last)
	for i := len(callers) - 1; i >= 0; i-- {
		fmt.Printf("%d: %s\n", len(callers)-i, callers[i])
	}

	return nil
}

var ErrCallerNotFound = fmt.Errorf("caller not found")

func getCallerName(skip int) (string, error) {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "", ErrCallerNotFound
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "", fmt.Errorf("unable to get function for PC")
	}
	fullName := f.Name()
	// Trim the path from the function name
	if lastSlash := strings.LastIndex(fullName, "/"); lastSlash >= 0 {
		fullName = fullName[lastSlash+1:]
	}
	return fullName, nil
}
