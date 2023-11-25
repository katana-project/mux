package mux

import (
	"github.com/katana-project/ffmpeg/avutil"
)

type Frame struct {
	frame *avutil.Frame
}

func NewFrame() *Frame {
	frm := &avutil.Frame{}
	if frm.Alloc(); frm.Null() {
		panic("memory allocation failed")
	}

	return &Frame{frame: frm}
}

func (f *Frame) Close() error {
	f.frame.Free()
	return nil
}

func (f *Frame) Clear() error {
	f.frame.Unref()
	return nil
}
