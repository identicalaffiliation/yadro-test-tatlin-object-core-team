package writer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (w *BucketWriter) WriteBucket(ctx context.Context, bucketIndex int, in <-chan string) {
	file, err := os.Create(filepath.Join(w.dirPattern, fmt.Sprintf("bucket_%d.txt", bucketIndex)))
	if err != nil {
		log.Printf("create bucket file: %v\n", err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close file: %v", err)
		}
	}()

	writer := bufio.NewWriter(file)
	linesBatch := make([]string, 0, w.batchSize)

	defer func() {
		flushLinesBatch(writer, linesBatch)
		_ = writer.Flush()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-in:
			if !ok {
				return
			}

			linesBatch = append(linesBatch, line)
			if len(linesBatch) >= w.batchSize {
				linesBatch = flushLinesBatch(writer, linesBatch)
			}
		}
	}
}

func flushLinesBatch(w *bufio.Writer, batch []string) []string {
	if len(batch) == 0 {
		return batch
	}

	flush(w, batch)
	return batch[:0]
}

func flush(w *bufio.Writer, batch []string) {
	for _, line := range batch {
		_, _ = w.WriteString(line)
		_ = w.WriteByte('\n')
	}
}
