package main

import (
	"context"
	"net"
	"sync"
	"time"
)

// RunUDPNoise runs the UDP noise test with the given configuration
func RunUDPNoise(ctx context.Context, config UDPConfig) error {
	if err := ValidateUDPConfig(config); err != nil {
		return err
	}

	var wg sync.WaitGroup
	noiseData := []byte("noise")

	// Spawn worker goroutines, each with its own UDP connection
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Each worker creates its own UDP connection
			conn, err := net.Dial("udp", config.Target)
			if err != nil {
				// If connection fails, worker exits silently
				// This allows other workers to continue
				return
			}
			defer conn.Close()

			// If rate is specified, use ticker for rate limiting
			if config.Rate > 0 {
				// Rate is per worker, so divide by number of workers
				workerRate := config.Rate / float64(config.Workers)
				interval := time.Duration(float64(time.Second) / workerRate)
				ticker := time.NewTicker(interval)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
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
						return
					default:
						if _, err := conn.Write(noiseData); err != nil {
							// Continue even if write fails (target might not be listening)
							continue
						}
					}
				}
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()
	return nil
}
