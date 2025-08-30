package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("🚀 Go CRUD Unit Test Runner")
	fmt.Println("==============================")

	err := os.Setenv("CGO_ENABLED", "0")
	if err != nil {
		return
	}

	// List of test packages to run
	testPackages := []string{
		"./booking_test",
		"./invertory_test",
		"./invoice_test",
		"./user_test",
		"./vehicle_test",
		"./servicehistory_test",
		"./main_test",
	}

	totalTests := 0
	passedTests := 0

	for _, testPkg := range testPackages {
		fmt.Printf("\n📋 Running tests in %s...\n", testPkg)

		cmd := exec.Command("go", "test", "-v", testPkg)
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("❌ %s: FAILED\n", testPkg)
			fmt.Printf("Error: %s\n", string(output))
		} else {
			fmt.Printf("✅ %s: PASSED\n", testPkg)
			passedTests++
		}
		totalTests++
	}

	fmt.Printf("\n📊 Test Summary:\n")
	fmt.Printf("Total test packages: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", totalTests-passedTests)

	if passedTests == totalTests {
		fmt.Println("🎉 All tests passed!")
	} else {
		fmt.Println("⚠️  Some tests failed. Check output above.")
		os.Exit(1)
	}
}
