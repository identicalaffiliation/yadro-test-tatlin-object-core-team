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

	if !strings.HasSuffix(strings.TrimSpace(filename), ".txt") {
		log.Printf("file is not a txt format\n")
		os.Exit(1)
	}

	out := make(chan string, cfg.ChannelBuff)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileReader := reader.NewFileReader(filename)
	if err := fileReader.ParseFile(ctx, out); err != nil {
		log.Printf("parse file: %v", err)
		os.Exit(1)
	}

}
