package main

import (
	"context"
	"crypto/sha256"
	"sync"
)

// RunCryptoLoad runs the crypto load test with the given configuration
func RunCryptoLoad(ctx context.Context, config CryptoConfig) error {
	if err := ValidateCryptoConfig(config); err != nil {
		return err
	}

	var wg sync.WaitGroup
	data := make([]byte, config.DataSize)

	// Spawn worker goroutines
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Perform SHA256 operation to simulate cryptographic workload
					h := sha256.Sum256(data)
					_ = h // Use the result to prevent optimization
				}
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()
	return nil
}
