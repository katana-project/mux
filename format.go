package mux

import (
	"fmt"
	"github.com/katana-project/ffmpeg"
	"github.com/katana-project/ffmpeg/avcodec"
	"github.com/katana-project/ffmpeg/avformat"
	"path/filepath"
	"slices"
	"strings"
)

type Format interface {
	Name() string
	LongName() string
	MIME() []string
	Extensions() []string
	Input() bool
}

type Muxer struct {
	fmt *avformat.OutputFormat
}

func AvailableMuxers() []*Muxer {
	var (
		iterState = &ffmpeg.IterationState{}
		muxs      []*Muxer
	)
	for {
		of := avformat.MuxerIterate(iterState)
		if of == nil {
			break
		}

		muxs = append(muxs, &Muxer{fmt: of})
	}

	return muxs
}

func FindMuxer(name, ext, mime string) *Muxer {
	if !strings.ContainsRune(ext, '.') {
		ext = "." + ext
	}

	result := avformat.GuessFormat(name, ext, mime)
	if result == nil {
		return nil
	}

	return &Muxer{fmt: result}
}

func MustFindMuxer(name, ext, mime string) *Muxer {
	result := FindMuxer(name, ext, mime)
	if result == nil {
		panic(fmt.Sprintf("could not find muxer with name: %s, ext: %s, mime: %s", name, ext, mime))
	}

	return result
}

func (m *Muxer) Name() string {
	return m.fmt.Name()
}

func (m *Muxer) LongName() string {
	return m.fmt.LongName()
}

func (m *Muxer) MIME() []string {
	mimes := m.fmt.MIMEType()
	if mimes == "" {
		return nil
	}

	return strings.Split(mimes, ",")
}

func (m *Muxer) Extensions() []string {
	exts := m.fmt.Extensions()
	if exts == "" {
		return nil
	}

	return strings.Split(exts, ",")
}

func (m *Muxer) Input() bool {
	return false
}

func (m *Muxer) SupportsCodec(codec *Codec) bool {
	return m.fmt.QueryCodec(codec.id, avcodec.ComplianceNormal) > 0
}

type Demuxer struct {
	fmt *avformat.InputFormat
}

func AvailableDemuxers() []*Demuxer {
	var (
		iterState = &ffmpeg.IterationState{}
		demuxs    []*Demuxer
	)
	for {
		ifo := avformat.DemuxerIterate(iterState)
		if ifo == nil {
			break
		}

		demuxs = append(demuxs, &Demuxer{fmt: ifo})
	}

	return demuxs
}

func FindDemuxer(name, ext, mime string) *Demuxer {
	if name == "" && ext == "" && mime == "" {
		return nil // FAST PATH: not searching for anything
	}
	if strings.ContainsRune(ext, '.') {
		ext = filepath.Ext(ext) // TODO: use av_match_ext
	}

	name = strings.ToLower(name)
	for _, demux := range AvailableDemuxers() {
		if (name != "" && name == strings.ToLower(demux.Name())) ||
			(ext != "" && slices.Contains(demux.Extensions(), ext)) ||
			(mime != "" && slices.Contains(demux.MIME(), mime)) {
			return demux
		}
	}

	return nil
}

func MustFindDemuxer(name, ext, mime string) *Demuxer {
	result := FindDemuxer(name, ext, mime)
	if result == nil {
		panic(fmt.Sprintf("could not find demuxer with name: %s, ext: %s, mime: %s", name, ext, mime))
	}

	return result
}

func (d *Demuxer) Name() string {
	return d.fmt.Name()
}

func (d *Demuxer) LongName() string {
	return d.fmt.LongName()
}

func (d *Demuxer) MIME() []string {
	mimes := d.fmt.MIMEType()
	if mimes == "" {
		return nil
	}

	return strings.Split(mimes, ",")
}

func (d *Demuxer) Extensions() []string {
	exts := d.fmt.Extensions()
	if exts == "" {
		return nil
	}

	return strings.Split(exts, ",")
}

func (d *Demuxer) Input() bool {
	return true
}
