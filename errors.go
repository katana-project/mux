package mux

import (
	"errors"
	"fmt"
	"github.com/katana-project/ffmpeg/avutil"
	"io"
)

var (
	ErrAgain = errors.New("resource temporarily unavailable")

	ErrNotOpen = errors.New("encoder/decoder is not open for I/O")
)

type ErrAV struct {
	Code int
}

func (eav *ErrAV) Error() string {
	return fmt.Sprintf("FFmpeg error %d (%s)", eav.Code, avutil.Err2Str(eav.Code))
}

func wrapErrorCode(code int) error {
	if code >= 0 {
		return nil
	}

	switch code {
	case avutil.ErrorENoMem:
		panic("memory allocation failed")
	case avutil.ErrorEOF:
		return io.EOF
	case avutil.ErrorEAgain:
		return ErrAgain
	}
	return &ErrAV{Code: code}
}
