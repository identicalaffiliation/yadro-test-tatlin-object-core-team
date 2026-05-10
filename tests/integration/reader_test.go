package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/reader"
)

const (
	testFilename = "test.txt"
	testContent  = "Влад\nВлад\nРома\nМаша\nРома"
	chanBuff     = 5
	filePerm     = 0666
)

func TestReader_ParseFile(t *testing.T) {
	tempDir := t.TempDir()
	filepath := filepath.Join(tempDir, testFilename)
	err := os.WriteFile(filepath, []byte(testContent), filePerm)
	assert.NoError(t, err)

	t.Run("success reading", func(t *testing.T) {
		fr := reader.NewFileReader(filepath)
		out := make(chan string, chanBuff)
		ctx := context.Background()

		err := fr.ParseFile(ctx, out)
		close(out)
		assert.NoError(t, err)

		var lines []string
		for line := range out {
			lines = append(lines, line)
		}

		expected := []string{"Влад", "Влад", "Рома", "Маша", "Рома"}
		assert.SliceEqualT(t, expected, lines)
	})

	t.Run("file not found", func(t *testing.T) {
		fr := reader.NewFileReader("fake.txt")
		out := make(chan string, chanBuff)
		ctx := context.Background()
		err := fr.ParseFile(ctx, out)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "open file")
	})
}
