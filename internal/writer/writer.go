package writer

import (
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/internal/config"
	"github.com/identicalaffiliation/yadro-test-tatlin-object-core-team/pkg/hasher"
)

type Hasher interface {
	HashLine32(line string) uint32
}

type BucketWriter struct {
	hasher           Hasher
	bucketCount      int
	batchSize        int
	bucketChanBuffer int
	dirPattern       string
}

func NewBucketWriter(cfg *config.CLIConfig, dir string) *BucketWriter {
	return &BucketWriter{
		bucketCount:      cfg.BucketCount,
		batchSize:        cfg.BatchSize,
		bucketChanBuffer: cfg.ChannelBuff,
		hasher:           hasher.NewHasher(),
		dirPattern:       dir,
	}
}
