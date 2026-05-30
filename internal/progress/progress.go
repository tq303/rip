package progress

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

func Bar(size int64) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		size,
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionThrottle(100*time.Millisecond),
	)
}
