package mux

import (
	"fmt"
	"github.com/katana-project/ffmpeg/avcodec"
)

type Codec struct {
	enc, dec *avcodec.Codec
}

func NewCodec(id avcodec.CodecID) *Codec {
	return &Codec{enc: avcodec.FindEncoder(id), dec: avcodec.FindDecoder(id)}
}

func FindCodec(name string) *Codec {
	var (
		enc = avcodec.FindEncoderByName(name)
		dec = avcodec.FindDecoderByName(name)
	)
	if enc == nil && dec == nil {
		return nil
	}

	return &Codec{enc: enc, dec: dec}
}

func MustFindCodec(name string) *Codec {
	result := FindCodec(name)
	if result == nil {
		panic(fmt.Sprintf("could not find codec with name: %s", name))
	}

	return result
}

func (c *Codec) NewEncoder() *CodecIO {
	if c.enc == nil {
		return nil
	}

	return c.newIO(c.enc, true)
}

func (c *Codec) NewDecoder() *CodecIO {
	if c.dec == nil {
		return nil
	}

	return c.newIO(c.dec, false)
}

func (c *Codec) newIO(codec *avcodec.Codec, encoder bool) *CodecIO {
	cio := &CodecIO{
		Codec:    c,
		codecCtx: &avcodec.CodecContext{},
		encoder:  encoder,
	}
	if cio.codecCtx.AllocContext3(codec); cio.codecCtx.Null() {
		panic("memory allocation failed")
	}

	return cio
}

type CodecIO struct {
	*Codec

	codecCtx      *avcodec.CodecContext
	encoder, open bool
}

func (cio *CodecIO) Open() error {
	codec := cio.dec
	if cio.encoder {
		codec = cio.enc
	}
	if err := wrapErrorCode(cio.codecCtx.Open2(codec, nil)); err != nil {
		return err
	}

	cio.open = true
	return nil
}

func (cio *CodecIO) Close() error {
	defer cio.codecCtx.FreeContext()
	if cio.open {
		cio.open = false

		if err := wrapErrorCode(cio.codecCtx.Close()); err != nil {
			return err
		}
	}

	return nil
}

func (cio *CodecIO) Encoder() bool {
	return cio.encoder
}

func (cio *CodecIO) ReadPacket(pkt *Packet) error {
	if !cio.open {
		return ErrNotOpen
	}

	return wrapErrorCode(cio.codecCtx.ReceivePacket(pkt.packet))
}

func (cio *CodecIO) WritePacket(pkt *Packet) error {
	if !cio.open {
		return ErrNotOpen
	}

	return wrapErrorCode(cio.codecCtx.SendPacket(pkt.packet))
}

func (cio *CodecIO) ReadFrame(frm *Frame) error {
	if !cio.open {
		return ErrNotOpen
	}

	return wrapErrorCode(cio.codecCtx.ReceiveFrame(frm.frame))
}

func (cio *CodecIO) WriteFrame(frm *Frame) error {
	if !cio.open {
		return ErrNotOpen
	}

	return wrapErrorCode(cio.codecCtx.SendFrame(frm.frame))
}

func (cio *CodecIO) CopyParameters(dst *CodecIO) error {
	dst.codecCtx.SetTimeBase(cio.codecCtx.Framerate().Inv())
	dst.codecCtx.SetPixFmt(cio.codecCtx.PixFmt())
	dst.codecCtx.SetWidth(cio.codecCtx.Width())
	dst.codecCtx.SetHeight(cio.codecCtx.Height())

	return nil
}

func (cio *CodecIO) CopyCodecParameters(dst *Stream) error {
	return wrapErrorCode(cio.codecCtx.ParametersFromContext(dst.stream.CodecPar()))
}

func (cio *CodecIO) Flush() error {
	if !cio.open {
		return ErrNotOpen
	}
	if cio.encoder {
		return wrapErrorCode(cio.codecCtx.SendFrame(nil))
	}

	pkt := &avcodec.Packet{}
	if pkt.Alloc(); pkt.Null() {
		panic("memory allocation failed")
	}
	defer pkt.Free()

	return wrapErrorCode(cio.codecCtx.SendPacket(pkt))
}
