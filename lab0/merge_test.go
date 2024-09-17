package lab0_test

import (
	"context"
	"testing"

	"cs426.cloud/lab0"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func chanToSlice[T any](ch chan T) []T {
	vals := make([]T, 0)
	for item := range ch {
		vals = append(vals, item)
	}
	return vals
}

type mergeFunc = func(chan string, chan string, chan string)

func runMergeTest(t *testing.T, merge mergeFunc) {
	t.Run("empty channels", func(t *testing.T) {
		a := make(chan string)
		b := make(chan string)
		out := make(chan string)
		close(a)
		close(b)

		merge(a, b, out)
		// If your lab0 hangs here, make sure you are closing your channels!
		require.Empty(t, chanToSlice(out))
	})


}

func TestMergeChannels(t *testing.T) {
	runMergeTest(t, func(a, b, out chan string) {
		lab0.MergeChannels(a, b, out)
	})

	// NEW TEST
	t.Run("alternating sends", func(t *testing.T) {
		a := make(chan string)
		b := make(chan string)
		out := make(chan string, 4)

		go func() {
			a <- "a1"
			b <- "b1"
			a <- "a2"
			b <- "b2"
			close(a)
			close(b)
		}()
		

		lab0.MergeChannels(a, b, out)
		// If your lab0 hangs here, make sure you are closing your channels!
		require.Equal(t, []string{"a1", "b1", "a2", "b2"}, chanToSlice(out))
	})


	// NEW TEST
	t.Run("one closed channel", func(t *testing.T) {
		a := make(chan string)
		b := make(chan string)
		out := make(chan string, 4)
	

		go func() {
			a <- "a1"
			close(a)
			b <- "b1"
			b <- "b2"
			b <- "b3"
			close(b)
		}()
		

		lab0.MergeChannels(a, b, out)
		// If your lab0 hangs here, make sure you are closing your channels!
		require.Equal(t, []string{"a1", "b1", "b2", "b3"}, chanToSlice(out))
	})
}

func TestMergeOrCancel(t *testing.T) {
	runMergeTest(t, func(a, b, out chan string) {
		_ = lab0.MergeChannelsOrCancel(context.Background(), a, b, out)
	})

	t.Run("already canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		a := make(chan string, 1)
		b := make(chan string, 1)
		out := make(chan string, 10)

		eg, _ := errgroup.WithContext(context.Background())
		eg.Go(func() error {
			return lab0.MergeChannelsOrCancel(ctx, a, b, out)
		})
		err := eg.Wait()
		a <- "a"
		b <- "b"

		require.Error(t, err)
		require.Equal(t, []string{}, chanToSlice(out))
	})

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		a := make(chan string)
		b := make(chan string)
		out := make(chan string, 10)

		eg, _ := errgroup.WithContext(context.Background())
		eg.Go(func() error {
			return lab0.MergeChannelsOrCancel(ctx, a, b, out)
		})
		a <- "a"
		b <- "b"
		cancel()

		err := eg.Wait()
		require.Error(t, err)
		require.Equal(t, []string{"a", "b"}, chanToSlice(out))
	})

	// NEW TEST
	//checks that canceling the context after the merging does not affect the final output
	t.Run("cancel after close", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		a := make(chan string, 2)
		b := make(chan string, 2)
		out := make(chan string, 10)

		a <- "a1"
		a <- "a2"
		b <- "b1"
		b <- "b2"
		close(a)
		close(b)

		eg, _ := errgroup.WithContext(context.Background())
		eg.Go(func() error {
			return lab0.MergeChannelsOrCancel(ctx, a, b, out)
		})

		err := eg.Wait()
		
		cancel()


		require.NoError(t, err)
		require.ElementsMatch(t, []string{"a1", "a2", "b1", "b2"}, chanToSlice(out))
	})
}

type channelFetcher struct {
	ch chan string
}

func newChannelFetcher(ch chan string) *channelFetcher {
	return &channelFetcher{ch: ch}
}

func (f *channelFetcher) Fetch() (string, bool) {
	v, ok := <-f.ch
	return v, ok
}

func TestMergeFetches(t *testing.T) {
	runMergeTest(t, func(a, b, out chan string) {
		lab0.MergeFetches(newChannelFetcher(a), newChannelFetcher(b), out)
	})
}

type emptyFetcher struct{}

func (f *emptyFetcher) Fetch() (string, bool) {
    return "", false // Always returns false, indicating no data
}

func TestMergeFetchesAdditional(t *testing.T) {
	// TODO: add your extra tests here
	runMergeTest(t, func(a, b, out chan string) {
		lab0.MergeFetches(newChannelFetcher(a), newChannelFetcher(b), out)
	})

	t.Run("one closed channel", func(t *testing.T) {
		a := &emptyFetcher{}
		b := &emptyFetcher{}
		out := make(chan string)


		lab0.MergeFetches(a, b, out)
		// If your lab0 hangs here, make sure you are closing your channels!
		require.Equal(t, []string{"a1", "b1", "b2", "b3"}, chanToSlice(out))
	})

}
