package progress

import "github.com/schollz/progressbar/v3"

func Bar(label string, size int64) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		size,
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription(label),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
	)
}
