package writer

import (
	"context"
	"sync"
)

func (w *BucketWriter) Run(ctx context.Context, in <-chan string) error {
	bucketChanns := w.initChanns()
	wg := &sync.WaitGroup{}

	w.startWriters(ctx, bucketChanns, wg)

	defer func() {
		for _, ch := range bucketChanns {
			close(ch)
		}
		wg.Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-in:
			if !ok {
				return nil
			}

			bucketIndex := int(w.hasher.HashLine32(line)) % w.bucketCount
			select {
			case bucketChanns[bucketIndex] <- line:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
