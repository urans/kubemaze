package recipe

import (
	"log/slog"
	"sync"
)

// SlowStartBatch tries to call the provided function a total of 'count' times,
// starting slow to check for errors, then speeding up if calls succeed.
//
// It groups the calls into batches, starting with a group of initialBatchSize.
// Within each batch, it may call the function multiple times concurrently.
//
// If a whole batch succeeds, the next batch may get exponentially larger.
// If there are any failures in a batch, all remaining batches are skipped
// after waiting for the current batch to complete.
//
// It returns the number of successful calls to the function.
func SlowStartBatch(count int, initialBatchSize int, fn func() error) (int, error) {
	successes, remaining := 0, count
	for batchSize := min(remaining, initialBatchSize); batchSize > 0; batchSize = min(2*batchSize, remaining) {
		wg := &sync.WaitGroup{}
		errs := make(chan error, batchSize)
		for i := 0; i < batchSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := fn(); err != nil {
					errs <- err
				}
			}()
		}
		wg.Wait()

		curSuccesses := batchSize - len(errs)
		successes += curSuccesses
		if len(errs) > 0 {
			return successes, <-errs
		}
		remaining -= batchSize
		slog.Info("SlowStartBatch", "batchSize", batchSize, "successes", successes, "remaining", remaining)
	}
	return successes, nil
}
