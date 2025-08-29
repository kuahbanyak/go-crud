package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("🚀 Go CRUD Unit Test Runner")
	fmt.Println("==============================")

	// Set environment to disable CGO
	os.Setenv("CGO_ENABLED", "0")

	// List of test files to run individually
	testFiles := []string{
		"user_test.go",
		"vehicle_test.go",
		"booking_test.go",
		"inventory_test.go",
		"invoice_test.go",
		"servicehistory_test.go",
	}

	totalTests := 0
	passedTests := 0

	for _, testFile := range testFiles {
		fmt.Printf("\n📋 Running %s...\n", testFile)

		cmd := exec.Command("go", "test", "-v", fmt.Sprintf("./test/%s", testFile), "./test/main_test.go")
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("❌ %s: FAILED\n", testFile)
			fmt.Printf("Error: %s\n", string(output))
		} else {
			fmt.Printf("✅ %s: PASSED\n", testFile)
			passedTests++
		}
		totalTests++
	}

	fmt.Printf("\n📊 Test Summary:\n")
	fmt.Printf("Total test files: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", totalTests-passedTests)

	if passedTests == totalTests {
		fmt.Println("🎉 All tests passed!")
	} else {
		fmt.Println("⚠️  Some tests failed. Check output above.")
		os.Exit(1)
	}
}
