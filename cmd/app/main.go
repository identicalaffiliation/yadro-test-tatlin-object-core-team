package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/config"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/reader"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/reducer"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/writer"
)

const (
	TEMP_DIR_NAME = "./temp"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to cli config file")
	flag.Parse()

	cfg, err := config.LoadCLIConfig(configPath)
	if err != nil {
		log.Printf("load config: %v", err)
		os.Exit(1)
	}

	inReader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a filename: ")
	filename, err := inReader.ReadString('\n')
	if err != nil {
		log.Printf("read filename: %v\n", err)
		os.Exit(1)
	}

	filename = strings.TrimSpace(filename)

	if !strings.HasSuffix(filename, ".txt") {
		log.Printf("file is not a txt format\n")
		os.Exit(1)
	}

	chann := make(chan string, cfg.ChannelBuff)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fileReader := reader.NewFileReader(filename)
		if err := fileReader.ParseFile(ctx, chann); err != nil {
			log.Fatal(err)
		}
		close(chann)
	}()

	if err := os.MkdirAll(TEMP_DIR_NAME, 0755); err != nil {
		log.Printf("create base temp dir: %v", err)
		os.Exit(1)
	}

	tempDir, err := os.MkdirTemp(TEMP_DIR_NAME, "")
	if err != nil {
		log.Printf("make temp dir: %v", err)
		os.Exit(1)
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("remove all files from temp dir: %v", err)
		}

		if err := os.Remove(TEMP_DIR_NAME); err != nil {
			log.Printf("remove dir: %v", err)
		}
	}()

	writer := writer.NewBucketWriter(cfg, tempDir)
	if err := writer.Run(ctx, chann); err != nil {
		log.Fatal(err)
	}

	reducer := reducer.NewReducer(tempDir)
	result, err := reducer.GetInfos()
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range result {
		fmt.Printf("%s:%d\n", key, value)
	}
}
