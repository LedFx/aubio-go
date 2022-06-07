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
SpecDesc
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
	mat *MatrixBuffer
}

func NewFilterBank(filters uint, win_s uint) *FilterBank {
	fbo := C.new_aubio_filterbank(C.uint_t(filters), C.uint_t(win_s))
	return &FilterBank{
		o:   fbo,
		buf: NewSimpleBuffer(filters),
		mat: NewMatrixBufferFromFmat(C.aubio_filterbank_get_coeffs(fbo)),
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

func (fb *FilterBank) GetCoeffs() [][]float64 {
	return fb.mat.GetChannels()
}

func (fb *FilterBank) SetMelCoeffsSlaney(sample uint) {
	C.aubio_filterbank_set_mel_coeffs_slaney(fb.o, C.smpl_t(sample))
}

func (fb *FilterBank) SetCoeffs(coeffs [][]float64) {
	// not sure if we can modify fb.mat directly, so we'll make a new one and use that to update
	mb := &MatrixBuffer{
		Height: fb.mat.Height,
		Length: fb.mat.Length,
		mat:    C.new_fmat(C.uint_t(fb.mat.Height), C.uint_t(fb.mat.Length)),
	}
	mb.SetChannels(coeffs)
	C.aubio_filterbank_set_coeffs(fb.o, mb.mat)
}

// The coeffs will be normalized by the triangles area which results in an uneven melbank.
// Recommended you call NormaliseCoeffs after setting triangle bands.
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

func (fb *FilterBank) Coeffs() *MatrixBuffer {
	return fb.mat
}

// Normalize the filterbank triangles to a consistent height for an even melbank.
func (fb *FilterBank) NormalizeCoeffs() {
	for i := uint(0); i < fb.mat.Height; i++ {
		channel := fb.mat.GetChannel(i)
		// find the max of the channel
		var max float64
		for pos := range channel {
			if channel[pos] > max {
				max = channel[pos]
			}
		}
		if max == 0 {
			continue
		}
		// then normalise all the heights of the channel to the maximum
		for pos := range channel {
			channel[pos] /= max
		}
		fb.mat.SetChannel(i, channel)
	}
}

// mfcc

type MFCC struct {
	o      *C.aubio_mfcc_t
	coeffs *SimpleBuffer
}

func NewMFCC(bufSize, samplerate, n_coeffs, n_filters uint) (*MFCC, error) {
	mfcc, err := C.new_aubio_mfcc(
		C.uint_t(bufSize),
		C.uint_t(n_filters),
		C.uint_t(n_coeffs),
		C.uint_t(samplerate),
	)
	if err != nil {
		return nil, err
	}
	return &MFCC{
		o:      mfcc,
		coeffs: NewSimpleBuffer(bufSize)}, nil
}

func (mfcc *MFCC) Free() {
	if mfcc.o != nil {
		C.del_aubio_mfcc(mfcc.o)
		mfcc.o = nil
	}
	if mfcc.coeffs != nil {
		mfcc.coeffs.Free()
		mfcc.coeffs = nil
	}
}

func (mfcc *MFCC) Coeffs() *SimpleBuffer {
	return mfcc.coeffs
}

func (mfcc *MFCC) Do(in *ComplexBuffer) {
	if mfcc != nil || mfcc.o != nil {
		C.aubio_mfcc_do(mfcc.o, in.data, mfcc.coeffs.vec)
	} else {
		log.Println("Called Do on empty MFCC. Maybe you called Free previously?")
	}
}

func (mfcc *MFCC) SetScale(scale float64) {
	C.aubio_mfcc_set_scale(mfcc.o, C.smpl_t(scale))
}

func (mfcc *MFCC) GetScale() float64 {
	return float64(C.aubio_mfcc_get_scale(mfcc.o))
}

func (mfcc *MFCC) SetPower(power float64) {
	C.aubio_mfcc_set_power(mfcc.o, C.smpl_t(power))
}

func (mfcc *MFCC) GetPower() float64 {
	return float64(C.aubio_mfcc_get_power(mfcc.o))
}

func (mfcc *MFCC) SetMelCoeffsSlaney() {
	C.aubio_mfcc_set_mel_coeffs_slaney(mfcc.o)
}

func (mfcc *MFCC) SetMelCoeffsHTK(fmin uint, fmax uint) {
	C.aubio_mfcc_set_mel_coeffs_htk(mfcc.o, C.smpl_t(fmin), C.smpl_t(fmax))
}

func (mfcc *MFCC) SetMelCoeffs(fmin uint, fmax uint) {
	C.aubio_mfcc_set_mel_coeffs(mfcc.o, C.smpl_t(fmin), C.smpl_t(fmax))
}

// phasvoc

type PhaseVoc struct {
	o     *C.aubio_pvoc_t
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
