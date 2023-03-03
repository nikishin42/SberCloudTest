package playlist

import (
	"fmt"
	"time"
)

func parseDuration(duration time.Duration) string {
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	formattedDuration := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return formattedDuration
}

func NewPlayList() PlayList {
	pl := PlayList{
		RemoteChan: make(chan int),
	}
	return pl
}
