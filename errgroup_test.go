package errgroup_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/heppu/errgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrGroup_MultipleErrors(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	eg := &errgroup.ErrGroup{}
	eg.Go(func() error { return err1 })
	eg.Go(func() error { return err2 })
	err := eg.Wait()

	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
}

func TestErrGroup_Mixed(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	eg := &errgroup.ErrGroup{}
	eg.Go(func() error { return nil })
	eg.Go(func() error { return err1 })
	eg.Go(func() error { return nil })
	eg.Go(func() error { return err2 })
	eg.Go(func() error { return nil })
	err := eg.Wait()

	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
}

func TestErrGroup_NoError(t *testing.T) {
	eg := &errgroup.ErrGroup{}
	eg.Go(func() error { return nil })
	eg.Go(func() error { return nil })
	eg.Go(func() error { return nil })
	err := eg.Wait()
	require.NoError(t, err)
}

func TestErrGroup_NoTask(t *testing.T) {
	eg := &errgroup.ErrGroup{}
	err := eg.Wait()
	require.NoError(t, err)
}

// This example fetches several URLs concurrently,
// using a WaitGroup to block until all the fetches are complete.
func ExampleErrGroup() {
	urls := []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.example.com/",
	}

	// Create a new ErrGroup.
	eg := &errgroup.ErrGroup{}
	for _, url := range urls {
		url := url
		// Launch a goroutine to fetch the URL.
		eg.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			return resp.Body.Close()
		})
	}

	// Wait for all HTTP fetches to complete.
	err := eg.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
