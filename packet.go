package mux

import "github.com/katana-project/ffmpeg/avcodec"

type Packet struct {
	packet *avcodec.Packet
}

func NewPacket() *Packet {
	pkt := &avcodec.Packet{}
	if pkt.Alloc(); pkt.Null() {
		panic("memory allocation failed")
	}

	return &Packet{packet: pkt}
}

func (p *Packet) Close() error {
	p.packet.Free()
	return nil
}

func (p *Packet) StreamIndex() int {
	return p.packet.StreamIndex()
}

func (p *Packet) SetStreamIndex(index int) {
	p.packet.SetStreamIndex(index)
}

func (p *Packet) Pos() int {
	return int(p.packet.Pos())
}

func (p *Packet) SetPos(pos int) {
	p.packet.SetPos(int64(pos))
}

func (p *Packet) ResetPos() {
	p.packet.SetPos(-1)
}

func (p *Packet) Rescale(src, dst *Rational) {
	p.packet.RescaleTs(src.r, dst.r)
}

func (p *Packet) Clear() error {
	p.packet.Unref()
	return nil
}
