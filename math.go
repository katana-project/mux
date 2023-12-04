package mux

import "github.com/katana-project/ffmpeg/avutil"

type Rational struct {
	r *avutil.Rational
}

func NewRational(num, den int) *Rational {
	return &Rational{r: avutil.MakeRational(num, den)}
}

func (r *Rational) Num() int {
	return r.r.Num()
}

func (r *Rational) Den() int {
	return r.r.Den()
}
