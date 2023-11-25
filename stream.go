package mux

import (
	"github.com/katana-project/ffmpeg/avformat"
	"github.com/katana-project/ffmpeg/avutil"
)

type MediaType int

const (
	MediaTypeUnknown MediaType = iota
	MediaTypeAudio
	MediaTypeVideo
	MediaTypeSubtitle
)

type Stream struct {
	stream *avformat.Stream
}

func (s *Stream) Index() int {
	return s.stream.Index()
}

func (s *Stream) Type() MediaType {
	switch s.stream.CodecPar().CodecType() {
	case avutil.MediaTypeAudio:
		return MediaTypeAudio
	case avutil.MediaTypeVideo:
		return MediaTypeVideo
	case avutil.MediaTypeSubtitle:
		return MediaTypeSubtitle
	}

	return MediaTypeUnknown
}

func (s *Stream) TimeBase() *Rational {
	r := &Rational{}
	r.read(s.stream.TimeBase())

	return r
}

func (s *Stream) Codec() *Codec {
	return NewCodec(s.stream.CodecPar().CodecID())
}

func (s *Stream) Decoder() (*CodecIO, error) {
	c := s.Codec().NewDecoder()
	if err := s.CopyCodecParameters(c); err != nil {
		c.Close()
		return nil, err
	}
	if err := c.Open(); err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

func (s *Stream) CopyParameters(dst *Stream) error {
	return wrapErrorCode(s.stream.CodecPar().Copy(dst.stream.CodecPar()))
}

func (s *Stream) CopyCodecParameters(dst *CodecIO) error {
	return wrapErrorCode(dst.codecCtx.ParametersToContext(s.stream.CodecPar()))
}
