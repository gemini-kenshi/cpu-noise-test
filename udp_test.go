package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

// RunUDPNoise runs the UDP noise test with the given configuration
func RunUDPNoise(ctx context.Context, config UDPConfig) error {
	if err := ValidateUDPConfig(config); err != nil {
		return err
	}

	// Create UDP connection
	conn, err := net.Dial("udp", config.Target)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	defer conn.Close()

	noiseData := []byte("noise")

	// If rate is specified, use ticker for rate limiting
	if config.Rate > 0 {
		interval := time.Duration(float64(time.Second) / config.Rate)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				if _, err := conn.Write(noiseData); err != nil {
					// Continue even if write fails (target might not be listening)
					continue
				}
			}
		}
	} else {
		// Unlimited rate - send as fast as possible
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if _, err := conn.Write(noiseData); err != nil {
					// Continue even if write fails (target might not be listening)
					continue
				}
			}
		}
	}
}
