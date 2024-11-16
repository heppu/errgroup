// Package errgroup is simple wrapper around sync.WaitGroup adding error handling.
// By default it uses errors.Join to combine errors, but you can override it with your own function for custom error types.
package errgroup

import (
	"errors"
	"sync"
)

// Join allows overwriting the default errors.Join function.
// This can be useful for custom multi error formatting.
var Join = errors.Join

// ErrGroup is a goroutine group that waits for all goroutines to finish and collects errors.
type ErrGroup struct {
	wg    sync.WaitGroup
	mutex sync.RWMutex
	err   error
}

// Go runs the given function in a goroutine.
func (eg *ErrGroup) Go(f func() error) {
	eg.wg.Add(1)

	go func() {
		defer eg.wg.Done()
		if err := f(); err != nil {
			eg.mutex.Lock()
			eg.err = Join(eg.err, err)
			eg.mutex.Unlock()
		}
	}()
}

// Wait waits for all goroutines to finish and returns all errors that occurred.
func (eg *ErrGroup) Wait() error {
	eg.wg.Wait()
	eg.mutex.RLock()
	err := eg.err
	eg.mutex.RUnlock()
	return err
}
