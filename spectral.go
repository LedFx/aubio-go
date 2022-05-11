/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

/*
MISSING:
AWhitening
 - all
DCT
 - all
FFT
 - all
PhaseVoc
 - get win, get hop, set window
Filterbank
 - set and get coeffs (fmat type required)
MFCC
 - all
SpecDesc
 - all
TSS
 - all
*/

package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

import (
	"log"
)

// fft

// type FFT struct {
// 	o *C.aubio_fft_t
// 	buf *SimpleBuffer
// }

// filterbank
type FilterBank struct {
	o   *C.aubio_filterbank_t
	buf *SimpleBuffer
}

func NewFilterBank(filters uint, win_s uint) *FilterBank {
	return &FilterBank{
		o:   C.new_aubio_filterbank(C.uint_t(filters), C.uint_t(win_s)),
		buf: NewSimpleBuffer(filters),
	}
}

func (fb *FilterBank) Do(in *ComplexBuffer) {
	if fb.o != nil {
		C.aubio_filterbank_do(fb.o, in.data, fb.buf.vec)
	} else {
		log.Println("Called Do on empty FilterBank. Maybe you called Free previously?")
	}
}

func (fb *FilterBank) SetNorm(norm float64) {
	C.aubio_filterbank_set_norm(fb.o, C.smpl_t(norm))
}

func (fb *FilterBank) GetNorm() float64 {
	return float64(C.aubio_filterbank_get_norm(fb.o))
}

func (fb *FilterBank) SetPower(power float64) {
	C.aubio_filterbank_set_power(fb.o, C.smpl_t(power))
}

func (fb *FilterBank) GetPower() float64 {
	return float64(C.aubio_filterbank_get_power(fb.o))
}

func (fb *FilterBank) SetMelCoeffsSlaney(sample uint) {
	C.aubio_filterbank_set_mel_coeffs_slaney(fb.o, C.smpl_t(sample))
}

func (fb *FilterBank) SetTriangleBands(freqs *SimpleBuffer, sample uint) {
	C.aubio_filterbank_set_triangle_bands(fb.o, freqs.vec, C.smpl_t(sample))
}

func (fb *FilterBank) SetMelCoeffsHTK(sample uint, fmin uint, fmax uint) {
	C.aubio_filterbank_set_mel_coeffs_htk(fb.o, C.smpl_t(sample), C.smpl_t(fmin), C.smpl_t(fmax))
}

func (fb *FilterBank) SetMelCoeffs(sample uint, fmin uint, fmax uint) {
	C.aubio_filterbank_set_mel_coeffs(fb.o, C.smpl_t(sample), C.smpl_t(fmin), C.smpl_t(fmax))
}

func (fb *FilterBank) Buffer() *SimpleBuffer {
	return fb.buf
}

// mfcc

// phasvoc

type PhaseVoc struct {
	o     *C.aubio_pvoc_t
	buf   *SimpleBuffer
	grain *ComplexBuffer
}

func NewPhaseVoc(bufSize, fftLen uint) (*PhaseVoc, error) {
	pvoc, err := C.new_aubio_pvoc(C.uint_t(bufSize), C.uint_t(fftLen))
	if err != nil {
		return nil, err
	}
	return &PhaseVoc{
		o:     pvoc,
		grain: NewComplexBuffer(bufSize)}, nil
}

func (pv *PhaseVoc) Free() {
	if pv.o != nil {
		C.del_aubio_pvoc(pv.o)
		pv.o = nil
	}
	if pv.grain != nil {
		pv.grain.Free()
		pv.grain = nil
	}
}

func (pv *PhaseVoc) Grain() *ComplexBuffer {
	return pv.grain
}

func (pv *PhaseVoc) Do(in *SimpleBuffer) {
	if pv != nil || pv.o != nil {
		C.aubio_pvoc_do(pv.o, in.vec, pv.grain.data)
	} else {
		log.Println("Called Do on empty PhaseVoc. Maybe you called Free previously?")
	}
}

func (pv *PhaseVoc) ReverseDo(out *SimpleBuffer) {
	if pv.o != nil {
		C.aubio_pvoc_rdo(pv.o, pv.grain.data, out.vec)
	} else {
		log.Println("Called ReverseDo on empty PhaseVoc. Maybe you called Free previously?")
	}
}

// statistics

// tss

type TSS struct {
	o     *C.aubio_tss_t
	buf   *ComplexBuffer
	trans *ComplexBuffer
	stead *ComplexBuffer
}

func NewTSS(bufSize, fftLen uint) (*TSS, error) {
	tss, err := C.new_aubio_tss(C.uint_t(bufSize), C.uint_t(fftLen))
	if err != nil {
		return nil, err
	}
	return &TSS{
		o:     tss,
		trans: NewComplexBuffer(bufSize),
		stead: NewComplexBuffer(bufSize)}, nil
}

func (tss *TSS) Free() {
	if tss.o != nil {
		C.del_aubio_tss(tss.o)
		tss.o = nil
	}
	if tss.trans != nil {
		tss.trans.Free()
		tss.trans = nil
	}
	if tss.stead != nil {
		tss.stead.Free()
		tss.stead = nil
	}
}

func (tss *TSS) Trans() *ComplexBuffer {
	return tss.trans
}

func (tss *TSS) Stead() *ComplexBuffer {
	return tss.stead
}

func (tss *TSS) Do(in *ComplexBuffer) {
	if tss != nil || tss.o != nil {
		C.aubio_tss_do(tss.o, in.data, tss.trans.data, tss.stead.data)
	} else {
		log.Println("Called Do on empty TSS. Maybe you called Free previously?")
	}
}

func (tss *TSS) SetThreshold(thrs float64) {
	C.aubio_tss_set_threshold(tss.o, C.smpl_t(thrs))
}

func (tss *TSS) SetAlpha(alpha float64) {
	C.aubio_tss_set_alpha(tss.o, C.smpl_t(alpha))
}

func (tss *TSS) SetBeta(beta float64) {
	C.aubio_tss_set_beta(tss.o, C.smpl_t(beta))
}
