package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse flags
	mode := flag.String("mode", "", "Test mode: 'crypto' or 'udp' (required)")

	// Crypto mode flags
	dataSize := flag.Int("data-size", 1024, "Data size in bytes per SHA256 operation (crypto mode)")
	workers := flag.Int("workers", 1, "Number of concurrent worker goroutines (crypto mode)")

	// UDP mode flags
	target := flag.String("target", "127.0.0.1:1", "Target address:port (udp mode)")
	rate := flag.Float64("rate", 0, "Send rate in packets per second (udp mode, 0 = unlimited)")

	flag.Parse()

	// Validate mode
	if *mode == "" {
		fmt.Fprintf(os.Stderr, "Error: --mode is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *mode != "crypto" && *mode != "udp" {
		fmt.Fprintf(os.Stderr, "Error: --mode must be 'crypto' or 'udp'\n")
		os.Exit(1)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start goroutine to handle signals
	go func() {
		sig := <-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived signal: %v. Shutting down gracefully...\n", sig)
		cancel()
	}()

	// Run the appropriate test mode
	var err error
	switch *mode {
	case "crypto":
		config := CryptoConfig{
			DataSize: *dataSize,
			Workers:  *workers,
		}
		if err = ValidateCryptoConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid crypto configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Starting crypto load test (data-size=%d, workers=%d). Press Ctrl+C to stop.\n", config.DataSize, config.Workers)
		err = RunCryptoLoad(ctx, config)

	case "udp":
		config := UDPConfig{
			Target: *target,
			Rate:   *rate,
		}
		if err = ValidateUDPConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid UDP configuration: %v\n", err)
			os.Exit(1)
		}
		rateStr := "unlimited"
		if config.Rate > 0 {
			rateStr = fmt.Sprintf("%.0f pps", config.Rate)
		}
		fmt.Fprintf(os.Stderr, "Starting UDP noise test (target=%s, rate=%s). Press Ctrl+C to stop.\n", config.Target, rateStr)
		err = RunUDPNoise(ctx, config)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Test completed.\n")
}
