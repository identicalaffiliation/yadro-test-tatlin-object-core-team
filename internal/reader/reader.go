package reader

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

type FileReader struct {
	filename string
}

func NewFileReader(filename string) *FileReader {
	return &FileReader{filename: filename}
}

func (fr *FileReader) ParseFile(ctx context.Context, out chan<- string) error {
	file, err := os.Open(fr.filename)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- scanner.Text():
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan result: %w", err)
	}

	return nil
}
