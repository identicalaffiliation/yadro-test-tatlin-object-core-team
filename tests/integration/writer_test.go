package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/config"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/writer"
)

const (
	bucketCount = 5
	batchSize   = 10
)

func TestBucketWriter_StartWriters(t *testing.T) {
	bucketWriter := writer.NewBucketWriter(&config.CLIConfig{
		BucketCount: bucketCount,
	}, t.TempDir())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ins := make([]chan string, bucketCount)
	for i := 0; i < bucketCount; i++ {
		ins[i] = make(chan string, chanBuff)
	}

	wg := &sync.WaitGroup{}

	bucketWriter.StartWriters(ctx, ins, wg)

	for i := 0; i < bucketCount; i++ {
		ins[i] <- "test"
		close(ins[i])
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout expired")
	}
}

func TestBucketWriter_WriteBucket(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("success write bucket", func(t *testing.T) {
		bucketWriter := writer.NewBucketWriter(&config.CLIConfig{
			BatchSize:   batchSize,
			BucketCount: bucketCount,
		}, tmpDir)

		ctx := context.Background()

		in := make(chan string)
		bucketIndex := 1
		done := make(chan struct{})

		go func() {
			bucketWriter.WriteBucket(ctx, bucketIndex, in)
			close(done)
		}()

		lines := []string{"Влад", "Влад", "Володя"}
		for _, line := range lines {
			in <- line
		}

		close(in)
		<-done

		file1 := fmt.Sprintf("bucket_%d.txt", bucketIndex)
		filepath := filepath.Join(tmpDir, file1)

		content, err := os.ReadFile(filepath)
		assert.NoError(t, err)

		epected := "Влад\nВлад\nВолодя\n"
		assert.Equal(t, epected, string(content))
	})
}

func TestBucketWriter_Run(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("full success run", func(t *testing.T) {
		bucketWriter := writer.NewBucketWriter(&config.CLIConfig{
			ChannelBuff: chanBuff,
			BatchSize:   batchSize,
			BucketCount: bucketCount,
		}, tmpDir)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		in := make(chan string)

		errCh := make(chan error, 1)
		go func() {
			errCh <- bucketWriter.Run(ctx, in)
		}()

		in <- "test1"
		in <- "test2"
		close(in)

		err := <-errCh
		assert.NoError(t, err)

		files, err := os.ReadDir(tmpDir)
		assert.NoError(t, err)

		assert.Len(t, files, bucketCount)
	})

	t.Run("context cancel", func(t *testing.T) {
		bucketWriter := writer.NewBucketWriter(&config.CLIConfig{
			BucketCount: bucketCount,
			ChannelBuff: chanBuff,
			BatchSize:   batchSize,
		}, tmpDir)

		ctx, cancel := context.WithCancel(context.Background())
		in := make(chan string)

		go func() {
			cancel()
		}()

		err := bucketWriter.Run(ctx, in)
		assert.ErrorIs(t, err, context.Canceled)
	})
}
