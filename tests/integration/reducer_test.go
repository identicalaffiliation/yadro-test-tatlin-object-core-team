package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/reducer"
)

const (
	contentFirstFile   = "Влад\nВлад\nВася"
	contentSecondFile  = "Петя\nДаша"
	firstTestFilename  = "test1.txt"
	secondTestFilename = "test2.txt"
)

func TestReducer_GetInfos(t *testing.T) {
	tempDir := t.TempDir()
	file1 := filepath.Join(tempDir, firstTestFilename)
	file2 := filepath.Join(tempDir, secondTestFilename)
	err := os.WriteFile(file1, []byte(contentFirstFile), filePerm)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte(contentSecondFile), filePerm)
	assert.NoError(t, err)

	expected1 := map[string]int{
		"Влад": 2,
		"Вася": 1,
		"Петя": 1,
		"Даша": 1,
	}

	expected2 := map[string]int{}

	t.Run("success reduce", func(t *testing.T) {
		reducer := reducer.NewReducer(tempDir)
		result, err := reducer.GetInfos()
		assert.NoError(t, err)

		assert.MapEqualT(t, expected1, result)
	})

	t.Run("success empty reduce", func(t *testing.T) {
		tempDir = t.TempDir()
		reducer := reducer.NewReducer(tempDir)
		result, err := reducer.GetInfos()
		assert.NoError(t, err)

		assert.Len(t, result, 0)
		assert.MapEqualT(t, expected2, result)
	})
}
