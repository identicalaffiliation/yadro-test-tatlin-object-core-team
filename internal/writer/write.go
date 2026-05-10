package writer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (w *BucketWriter) writeBucket(_ context.Context, bucketIndex int, in <-chan string) {
	file, err := os.Create(filepath.Join(w.dirPattern, fmt.Sprintf("bucket_%d.txt", bucketIndex)))
	if err != nil {
		log.Printf("create bucket file: %v\n", err)
		return
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	linesBatch := make([]string, 0, w.batchSize)

	defer func() {
		flushLinesBatch(writer, linesBatch)
		_ = writer.Flush()
	}()

	for line := range in {
		linesBatch = append(linesBatch, line)

		if len(linesBatch) >= w.batchSize {
			linesBatch = flushLinesBatch(writer, linesBatch)
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
