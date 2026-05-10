package writer

func (w *BucketWriter) initChanns() []chan string {
	chans := make([]chan string, w.bucketCount)
	for i := 0; i < w.bucketCount; i++ {
		chans[i] = make(chan string, w.bucketChanBuffer)
	}

	return chans
}
