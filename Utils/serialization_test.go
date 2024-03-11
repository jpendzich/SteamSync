package utils

import "testing"

func TestDeserializeInt(test *testing.T) {
	input := make([]byte, 8)
	input[7] = 1
	output := DeserializeInt(input)
	if output.payload != 1 {
		test.Fatalf("Test Failed")
	}
}
