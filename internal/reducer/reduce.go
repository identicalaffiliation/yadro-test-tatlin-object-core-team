package reducer

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
)

func (r *Reducer) GetInfos() (map[string]int, error) {
	files, err := os.ReadDir(r.dirName)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	bucketFileInfos := make(map[string]int)
	for _, file := range files {
		bucketFileInfo, err := r.getInfo(file.Name())
		if err != nil {
			return nil, err
		}

		maps.Copy(bucketFileInfos, bucketFileInfo)
	}

	return bucketFileInfos, nil
}

func (r *Reducer) getInfo(filename string) (map[string]int, error) {
	fullPath := filepath.Join(r.dirName, filename)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open bucket file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close bucket file: %v", err)
		}
	}()

	bucketFileInfo := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bucketFileInfo[scanner.Text()]++
	}

	return bucketFileInfo, nil
}
