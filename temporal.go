/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

// Filter is a wrapper for the aubio_filter_t object.
type Filter struct {
	o   *C.aubio_filter_t
	buf *SimpleBuffer
}

// Constructs a Filter. Filters maintain their own working buffer
// which will get freed when the Filters Free method is called.
// The caller is responsible for calling Free on the constructed
// Filter or risk leaking memory.
func NewFilter(order, bufSize uint) (*Filter, error) {
	f, err := C.new_aubio_filter(C.uint_t(order))
	if f == nil {
		return nil, err
	}
	return &Filter{o: f, buf: NewSimpleBuffer(bufSize)}, nil
}

// Free frees up the memory allocatd by aubio for this Filter.
func (f *Filter) Free() {
	if f.o != nil {
		C.del_aubio_filter(f.o)
		f.o = nil
	}
	if f.buf != nil {
		f.buf.Free()
		f.buf = nil
	}
}

// Reset resets the memory for this Filter.
func (f *Filter) Reset() {
	if f.o != nil {
		C.aubio_filter_do_reset(f.o)
	}
}

// Buffer returns the output buffer for this Filter.
// The buffer is populated by calls to DoOutplace and is owned
// by the Filter object.
// Subsequent calls to DoOutplace may change the data contained in
// this buffer.
func (f *Filter) Buffer() *SimpleBuffer {
	return f.buf
}

// SetSamplerate sets the samplerate for this Filter.
func (f *Filter) SetSamplerate(rate uint) {
	if f.o != nil {
		C.aubio_filter_set_samplerate(f.o, C.uint_t(rate))
	}
}

// Do does an in-place filter on the input vector.
// The output buffer is not used.
func (f *Filter) Do(in *SimpleBuffer) {
	// Filter in-place
	if f.o != nil {
		C.aubio_filter_do(f.o, in.vec)
	}
}

// TODO(jwall): maybe the outplace filter should be a seperate type?
// DoOutPlace does filters the input vector into the Filter's output
// Buffer. Each call to this method will change the data contained
// in the output buffer. This buffer can be retrieved though the
// Buffer method.
func (f *Filter) DoOutplace(in *SimpleBuffer) {
	if f.o != nil {
		C.aubio_filter_do_outplace(f.o, in.vec, f.buf.vec)
	}
}

// DoFwdBack runs the aubio_filter_do_filtfilt function on this
// Filter.
func (f *Filter) DoFwdBack(in *SimpleBuffer, workBufSize uint) {
	if f.o != nil {
		tmp := NewSimpleBuffer(workBufSize)
		defer tmp.Free()
		C.aubio_filter_do_filtfilt(f.o, in.vec, tmp.vec)
	}
}

// Feedback returns the buffer containing the feedback coefficients.
func (f *Filter) Feedback() *LongSampleBuffer {
	if f.o != nil {
		return newLBufferFromVec(C.aubio_filter_get_feedback(f.o))
	}
	return nil
}

// Feedback returns the buffer containing the feedforward coefficients.
func (f *Filter) Feedforward() *LongSampleBuffer {
	if f.o != nil {
		return newLBufferFromVec(C.aubio_filter_get_feedforward(f.o))
	}
	return nil
}

// Order returns this Filters order.
func (f *Filter) Order() uint {
	if f.o != nil {
		return uint(C.aubio_filter_get_order(f.o))
	}
	return 0
}

// Samplerate returns this Filters samplerate.
func (f *Filter) Samplerate() uint {
	if f.o != nil {
		return uint(C.aubio_filter_get_samplerate(f.o))
	}
	return 0
}

// A-Weighting

// Constructs an A-design Filter
// Samplerate should be one of 8000, 11025, 16000, 22050, 24000, 32000, 44100, 48000, 88200, 96000, and 192000 Hz
func NewFilterAWeighting(samplerate, bufSize uint) (*Filter, error) {
	f, err := C.new_aubio_filter_a_weighting(C.uint_t(samplerate))
	if f == nil {
		return nil, err
	}
	return &Filter{o: f, buf: NewSimpleBuffer(bufSize)}, nil
}

// Apply A-weighting to a filter
func (f *Filter) SetAWeighting(samplerate uint) {
	C.aubio_filter_set_a_weighting(f.o, C.uint_t(samplerate))
}

// C-Weighting

// Constructs an C-design Filter
// Samplerate should be one of 8000, 11025, 16000, 22050, 24000, 32000, 44100, 48000, 88200, 96000, and 192000 Hz
func NewFilterCWeighting(samplerate, bufSize uint) (*Filter, error) {
	f, err := C.new_aubio_filter_c_weighting(C.uint_t(samplerate))
	if f == nil {
		return nil, err
	}
	return &Filter{o: f, buf: NewSimpleBuffer(bufSize)}, nil
}

// Apply C-weighting to a filter
func (f *Filter) SetCWeighting(samplerate uint) {
	C.aubio_filter_set_c_weighting(f.o, C.uint_t(samplerate))
}

// Biquad

// Constructs a biquad Filter
func NewFilterBiquad(
	b0 float64,
	b1 float64,
	b2 float64,
	a0 float64,
	a1 float64,
	bufSize uint,
) (*Filter, error) {
	f, err := C.new_aubio_filter_biquad(
		C.lsmp_t(b0),
		C.lsmp_t(b1),
		C.lsmp_t(b2),
		C.lsmp_t(a0),
		C.lsmp_t(a1),
	)
	if f == nil {
		return nil, err
	}
	return &Filter{o: f, buf: NewSimpleBuffer(bufSize)}, nil
}

// Apply biquad to a filter
func (f *Filter) SetBiquad(
	b0 float64,
	b1 float64,
	b2 float64,
	a0 float64,
	a1 float64,
) {
	C.aubio_filter_set_biquad(
		f.o,
		C.lsmp_t(b0),
		C.lsmp_t(b1),
		C.lsmp_t(b2),
		C.lsmp_t(a0),
		C.lsmp_t(a1),
	)
}

// Resampler
// Ratio: output_sample_rate / input_sample_rate
// Quality: Resample quality (0 is best, 4 is fastest)
type Resampler struct {
	o   *C.aubio_resampler_t
	buf *SimpleBuffer
}

// bufSize is the size of the input buffer.
// output buffer is size bufSize * ratio
func NewResampler(ratio float64, quality uint, bufSize uint) (*Resampler, error) {
	r, err := C.new_aubio_resampler(C.smpl_t(ratio), C.uint_t(quality))
	if r == nil {
		return nil, err
	}
	return &Resampler{o: r, buf: NewSimpleBuffer(uint(float64(bufSize) * ratio))}, nil
}

func (r *Resampler) Free() {
	if r.o != nil {
		C.del_aubio_resampler(r.o)
		r.o = nil
	}
	if r.buf != nil {
		r.buf.Free()
		r.buf = nil
	}
}

func (r *Resampler) Buffer() *SimpleBuffer {
	return r.buf
}

func (r *Resampler) Do(in *SimpleBuffer) {
	if r.o != nil {
		C.aubio_resampler_do(r.o, in.vec, r.buf.vec)
	}
}
