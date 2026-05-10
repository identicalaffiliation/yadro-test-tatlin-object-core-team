package writer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

func (w *BucketWriter) writeBucket(ctx context.Context, bucketIndex int, in <-chan string) {
	file, err := os.Create(fmt.Sprintf("./temp/bucket_%d.txt", bucketIndex))
	if err != nil {
		log.Printf("create bucket file: %v\n", err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close file: %v\n", err)
		}
	}()

	writer := bufio.NewWriter(file)

	defer func() {
		if err := writer.Flush(); err != nil {
			log.Printf("flush writer: %v\n", err)
		}
	}()

	linesBatch := make([]string, 0, w.batchSize)
	for {
		select {
		case <-ctx.Done():
			flushLinesBatch(writer, linesBatch)
			return
		case line, ok := <-in:
			if !ok {
				flushLinesBatch(writer, linesBatch)
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
		w.WriteString(line) //nolint: errcheck
		w.WriteByte('\n')   //nolint: errcheck
	}
	w.Flush() //nolint: errcheck
}
