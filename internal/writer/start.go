package writer

import (
	"context"
	"sync"
)

func (w *BucketWriter) StartWriters(ctx context.Context, ins []chan string, wg *sync.WaitGroup) {
	for i := range ins {
		wg.Add(1)
		go func(ind int, in chan string) {
			defer wg.Done()

			w.WriteBucket(ctx, ind, in)
		}(i, ins[i])
	}
}
