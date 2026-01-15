package main

import (
	"context"
	"testing"
	"time"
)

func TestRunCryptoLoad_ValidConfig(t *testing.T) {
	ctx := context.Background()
	config := CryptoConfig{
		DataSize: 1024,
		Workers:  2,
	}

	// Create a context that will be cancelled after a short time
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err := RunCryptoLoad(ctx, config)
	if err != nil {
		t.Fatalf("RunCryptoLoad returned error: %v", err)
	}
}

func TestRunCryptoLoad_InvalidDataSize(t *testing.T) {
	ctx := context.Background()
	config := CryptoConfig{
		DataSize: 0,
		Workers:  1,
	}

	err := RunCryptoLoad(ctx, config)
	if err == nil {
		t.Fatal("RunCryptoLoad should return error for invalid data size")
	}
}

func TestRunCryptoLoad_InvalidWorkers(t *testing.T) {
	ctx := context.Background()
	config := CryptoConfig{
		DataSize: 1024,
		Workers:  0,
	}

	err := RunCryptoLoad(ctx, config)
	if err == nil {
		t.Fatal("RunCryptoLoad should return error for invalid workers")
	}
}

func TestRunCryptoLoad_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := CryptoConfig{
		DataSize: 1024,
		Workers:  3,
	}

	// Start the test in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- RunCryptoLoad(ctx, config)
	}()

	// Cancel after a short delay
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for the function to return
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("RunCryptoLoad returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("RunCryptoLoad did not respond to context cancellation")
	}
}

func TestRunCryptoLoad_MultipleWorkers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	config := CryptoConfig{
		DataSize: 512,
		Workers:  5,
	}

	start := time.Now()
	err := RunCryptoLoad(ctx, config)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("RunCryptoLoad returned error: %v", err)
	}

	// Verify it ran for at least some time (workers were active)
	if duration < 50*time.Millisecond {
		t.Logf("Warning: Test completed very quickly (%v), workers may not have started", duration)
	}
}

func TestRunCryptoLoad_DifferentDataSizes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	testCases := []struct {
		name     string
		dataSize int
		workers  int
	}{
		{"Small data", 64, 1},
		{"Medium data", 1024, 2},
		{"Large data", 4096, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := CryptoConfig{
				DataSize: tc.dataSize,
				Workers:  tc.workers,
			}

			err := RunCryptoLoad(ctx, config)
			if err != nil {
				t.Fatalf("RunCryptoLoad returned error: %v", err)
			}
		})
	}
}
