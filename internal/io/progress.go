package io

import (
	"io"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar interface {
	io.Writer
}

// NewProgressBar creates a default bytes-based progress bar
func NewProgressBar() ProgressBar {
	return progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription("transferring"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
	)
}
