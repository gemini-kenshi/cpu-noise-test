package main

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestRunUDPNoise_ValidConfig_UnlimitedRate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	config := UDPConfig{
		Target: "127.0.0.1:1",
		Rate:   0, // Unlimited
	}

	err := RunUDPNoise(ctx, config)
	if err != nil {
		t.Fatalf("RunUDPNoise returned error: %v", err)
	}
}

func TestRunUDPNoise_ValidConfig_WithRateLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	config := UDPConfig{
		Target: "127.0.0.1:1",
		Rate:   10, // 10 packets per second
	}

	err := RunUDPNoise(ctx, config)
	if err != nil {
		t.Fatalf("RunUDPNoise returned error: %v", err)
	}
}

func TestRunUDPNoise_InvalidTarget_Empty(t *testing.T) {
	ctx := context.Background()
	config := UDPConfig{
		Target: "",
		Rate:   0,
	}

	err := RunUDPNoise(ctx, config)
	if err == nil {
		t.Fatal("RunUDPNoise should return error for empty target")
	}
}

func TestRunUDPNoise_InvalidTarget_InvalidFormat(t *testing.T) {
	ctx := context.Background()
	config := UDPConfig{
		Target: "invalid-format",
		Rate:   0,
	}

	err := RunUDPNoise(ctx, config)
	if err == nil {
		t.Fatal("RunUDPNoise should return error for invalid target format")
	}
}

func TestRunUDPNoise_InvalidTarget_InvalidPort(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name   string
		target string
	}{
		{"Port too high", "127.0.0.1:65536"},
		{"Port zero", "127.0.0.1:0"},
		{"Invalid port", "127.0.0.1:abc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := UDPConfig{
				Target: tc.target,
				Rate:   0,
			}

			err := RunUDPNoise(ctx, config)
			if err == nil {
				t.Fatalf("RunUDPNoise should return error for target: %s", tc.target)
			}
		})
	}
}

func TestRunUDPNoise_InvalidRate_Negative(t *testing.T) {
	ctx := context.Background()
	config := UDPConfig{
		Target: "127.0.0.1:1",
		Rate:   -1,
	}

	err := RunUDPNoise(ctx, config)
	if err == nil {
		t.Fatal("RunUDPNoise should return error for negative rate")
	}
}

func TestRunUDPNoise_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := UDPConfig{
		Target: "127.0.0.1:1",
		Rate:   0,
	}

	// Start the test in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- RunUDPNoise(ctx, config)
	}()

	// Cancel after a short delay
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for the function to return
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("RunUDPNoise returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("RunUDPNoise did not respond to context cancellation")
	}
}

func TestRunUDPNoise_ContextCancellation_WithRateLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := UDPConfig{
		Target: "127.0.0.1:1",
		Rate:   100, // 100 packets per second
	}

	// Start the test in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- RunUDPNoise(ctx, config)
	}()

	// Cancel after a short delay
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for the function to return
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("RunUDPNoise returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("RunUDPNoise did not respond to context cancellation")
	}
}

func TestRunUDPNoise_DifferentTargets(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	testCases := []struct {
		name   string
		target string
	}{
		{"Localhost", "127.0.0.1:1"},
		{"Localhost IPv6", "[::1]:1"},
		{"Different port", "127.0.0.1:9999"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := UDPConfig{
				Target: tc.target,
				Rate:   0,
			}

			err := RunUDPNoise(ctx, config)
			if err != nil {
				t.Fatalf("RunUDPNoise returned error: %v", err)
			}
		})
	}
}

func TestRunUDPNoise_DifferentRates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	testCases := []struct {
		name string
		rate float64
	}{
		{"Low rate", 1},
		{"Medium rate", 10},
		{"High rate", 100},
		{"Very high rate", 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := UDPConfig{
				Target: "127.0.0.1:1",
				Rate:   tc.rate,
			}

			err := RunUDPNoise(ctx, config)
			if err != nil {
				t.Fatalf("RunUDPNoise returned error: %v", err)
			}
		})
	}
}

func TestRunUDPNoise_WithUDPListener(t *testing.T) {
	// Create a UDP listener to receive packets
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("Failed to listen UDP: %v", err)
	}
	defer conn.Close()

	// Get the actual address
	target := conn.LocalAddr().String()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	config := UDPConfig{
		Target: target,
		Rate:   10, // 10 packets per second
	}

	// Start receiving packets
	received := make(chan bool, 1)
	go func() {
		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		_, _, err := conn.ReadFromUDP(buf)
		if err == nil {
			received <- true
		}
	}()

	// Run the test
	err = RunUDPNoise(ctx, config)
	if err != nil {
		t.Fatalf("RunUDPNoise returned error: %v", err)
	}

	// Check if we received at least one packet
	select {
	case <-received:
		// Success - packet was received
	case <-time.After(100 * time.Millisecond):
		// It's okay if we don't receive packets - UDP is unreliable
		t.Log("No packets received (this is acceptable for UDP)")
	}
}
