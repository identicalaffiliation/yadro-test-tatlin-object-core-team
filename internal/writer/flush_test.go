package writer

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestFlush(t *testing.T) {
	t.Run("success flush", func(t *testing.T) {
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)
		batch := []string{"test1", "test2"}

		flushBatch := flushLinesBatch(writer, batch)
		err := writer.Flush()
		assert.NoError(t, err)

		expected := "test1\ntest2\n"
		assert.Equal(t, expected, buffer.String())
		assert.Len(t, flushBatch, 0)
		assert.Equal(t, cap(flushBatch), cap(batch))
	})

	t.Run("empty batch", func(t *testing.T) {
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)
		batch := []string{}

		flushBatch := flushLinesBatch(writer, batch)
		err := writer.Flush()
		assert.NoError(t, err)

		assert.Len(t, buffer.Bytes(), 0)
		assert.Len(t, flushBatch, 0)
	})
}
