package analysis

import (
	"context"
	"fmt"
)

// Stop stops the analyzer service gracefully.
func (a *analyzer) Stop(ctx context.Context) error {
	a.logger.Info("stopping LLM analyzer service")

	a.cancel() // cancel the context to stop receiving new feedbacks and wait for other goroutines to finish
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		a.logger.Info("analyzer service stopped gracefully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout waiting for analyzer to stop")
	}
}
