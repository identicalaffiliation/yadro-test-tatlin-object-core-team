package writer

import (
	"context"
	"sync"
)

func (w *BucketWriter) startWriters(ctx context.Context, ins []chan string, wg *sync.WaitGroup) {
	for i := range ins {
		wg.Add(1)
		go func(ind int, in chan string) {
			defer wg.Done()

			w.writeBucket(ctx, ind, in)
		}(i, ins[i])
	}
}
