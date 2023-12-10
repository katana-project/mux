package mux

import (
	"github.com/katana-project/ffmpeg/avcodec"
	"github.com/katana-project/ffmpeg/avformat"
	"github.com/katana-project/ffmpeg/avutil"
)

type IOContext struct {
	formatCtx *avformat.FormatContext
	outputFmt *avformat.OutputFormat
}

func NewInputContext(path string) (*IOContext, error) {
	fc := &avformat.FormatContext{}
	if err := wrapErrorCode(fc.OpenInput(path, nil, nil)); err != nil {
		return nil, err
	}
	if err := wrapErrorCode(fc.FindStreamInfo(nil)); err != nil {
		fc.CloseInput()
		return nil, err
	}

	return &IOContext{formatCtx: fc}, nil
}

func NewOutputContext(fmt *Muxer, path string) (*IOContext, error) {
	fc := &avformat.FormatContext{}
	if err := wrapErrorCode(fc.AllocOutputContext2(fmt.fmt, "", path)); err != nil {
		return nil, err
	}

	ofmt := fc.OFormat()
	if (ofmt.Flags() & avformat.FmtNoFile) == 0 {
		pb := fc.Pb()
		if err := wrapErrorCode(pb.Open(path, avformat.IOFlagWrite)); err != nil {
			fc.FreeContext()
			return nil, err
		}

		fc.SetPb(pb)
	}

	return &IOContext{formatCtx: fc, outputFmt: ofmt}, nil
}

func (ioc *IOContext) Close() error {
	if ioc.outputFmt != nil {
		defer ioc.formatCtx.FreeContext()
		if (ioc.outputFmt.Flags() & avformat.FmtNoFile) == 0 {
			pb := ioc.formatCtx.Pb()
			if err := wrapErrorCode(pb.CloseP()); err != nil {
				return err
			}

			ioc.formatCtx.SetPb(pb)
		}
	} else {
		ioc.formatCtx.CloseInput()
	}

	return nil
}

func (ioc *IOContext) Metadata() map[string]string {
	var (
		metadata  = ioc.formatCtx.Metadata()
		metadata0 = make(map[string]string)
		tag       *avutil.DictionaryEntry
	)
	for {
		tag = metadata.Iterate(tag)
		if tag == nil {
			break
		}

		metadata0[tag.Key()] = tag.Value()
	}

	return metadata0
}

func (ioc *IOContext) Streams() []*Stream {
	var (
		streams  = ioc.formatCtx.Streams()
		streams0 = make([]*Stream, len(streams))
	)
	for i, s := range streams {
		streams0[i] = &Stream{stream: s}
	}

	return streams0
}

func (ioc *IOContext) Stream(index int) *Stream {
	s := ioc.formatCtx.Stream(index)
	if s == nil {
		return nil
	}

	return &Stream{stream: s}
}

func (ioc *IOContext) NewStream(codec *Codec) *Stream {
	var codec0 *avcodec.Codec
	if codec != nil {
		codec0 = codec.enc
	}

	s := ioc.formatCtx.NewStream(codec0)
	if s == nil {
		panic("memory allocation failed")
	}

	return &Stream{stream: s}
}

func (ioc *IOContext) WriteHeader() error {
	return wrapErrorCode(ioc.formatCtx.WriteHeader(nil))
}

func (ioc *IOContext) ReadFrame(pkt *Packet) error {
	return wrapErrorCode(ioc.formatCtx.ReadFrame(pkt.packet))
}

func (ioc *IOContext) WriteFrame(pkt *Packet) error {
	return wrapErrorCode(ioc.formatCtx.InterleavedWriteFrame(pkt.packet))
}

func (ioc *IOContext) WriteEnd() error {
	return wrapErrorCode(ioc.formatCtx.WriteTrailer())
}
