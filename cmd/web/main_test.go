package main

import (
	"log"
	"testing"
)

func TestRun(t *testing.T) {
	// Test the run function to ensure it initializes the application correctly
	_, err := run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	log.Println("Run function executed successfully")
}
