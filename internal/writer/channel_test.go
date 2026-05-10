package writer

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestBucketWriter_InitChanns(t *testing.T) {
	testCases := []struct {
		Name          string
		BucketCount   int
		ChannelBuffer int
	}{
		{
			Name:          "success",
			BucketCount:   5,
			ChannelBuffer: 10,
		},
		{
			Name:          "unbuf chan",
			BucketCount:   5,
			ChannelBuffer: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bucketWriter := BucketWriter{
				bucketCount:      tc.BucketCount,
				bucketChanBuffer: tc.ChannelBuffer,
			}

			bucketChanns := bucketWriter.initChanns()
			assert.Len(t, bucketChanns, tc.BucketCount)

			for _, channel := range bucketChanns {
				assert.Equal(t, tc.ChannelBuffer, cap(channel))
			}
		})
	}
}
