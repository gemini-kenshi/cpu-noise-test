package main

import (
	"fmt"
	"net"
	"strconv"
)

// CryptoConfig holds configuration for crypto load test
type CryptoConfig struct {
	DataSize int // Data size in bytes per SHA256 operation
	Workers  int // Number of concurrent worker goroutines
}

// UDPConfig holds configuration for UDP noise test
type UDPConfig struct {
	Target  string  // Target address:port (e.g., "127.0.0.1:1")
	Rate    float64 // Send rate in packets per second (0 means unlimited)
	Workers int     // Number of concurrent worker goroutines
}

// ValidateCryptoConfig validates crypto configuration
func ValidateCryptoConfig(config CryptoConfig) error {
	if config.DataSize <= 0 {
		return fmt.Errorf("data-size must be greater than 0")
	}
	if config.Workers <= 0 {
		return fmt.Errorf("workers must be greater than 0")
	}
	return nil
}

// ValidateUDPConfig validates UDP configuration
func ValidateUDPConfig(config UDPConfig) error {
	if config.Target == "" {
		return fmt.Errorf("target cannot be empty")
	}

	// Validate address:port format
	host, port, err := net.SplitHostPort(config.Target)
	if err != nil {
		return fmt.Errorf("invalid target format: %w", err)
	}

	if host == "" {
		return fmt.Errorf("target host cannot be empty")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid port number: %w", err)
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if config.Rate < 0 {
		return fmt.Errorf("rate cannot be negative")
	}

	if config.Workers <= 0 {
		return fmt.Errorf("workers must be greater than 0")
	}

	return nil
}
