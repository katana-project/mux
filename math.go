package mux

import "github.com/katana-project/ffmpeg/avutil"

type Rational struct {
	Num, Den int
}

func (r *Rational) read(rational *avutil.Rational) {
	r.Num = rational.Num()
	r.Den = rational.Den()
}

func (r *Rational) write() *avutil.Rational {
	return avutil.MakeRational(r.Num, r.Den)
}
