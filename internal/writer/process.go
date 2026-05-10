package writer

import (
	"context"
	"sync"
)

func (w *BucketWriter) Run(ctx context.Context, in <-chan string) error {
	bucketChanns := w.initChanns()
	wg := &sync.WaitGroup{}

	w.StartWriters(ctx, bucketChanns, wg)

	go func() {
		defer func() {
			for _, ch := range bucketChanns {
				close(ch)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return

			case line, ok := <-in:
				if !ok {
					return
				}

				bucketIndex := int(w.hasher.HashLine32(line)) % w.bucketCount

				select {
				case bucketChanns[bucketIndex] <- line:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	wg.Wait()
	return ctx.Err()
}
