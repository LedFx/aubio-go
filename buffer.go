package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"
import "unsafe"

// SimpleBuffer is a wrapper for the aubio fvec_t type. It is used
// as the buffer for processing audio data in an aubio pipeline.
// It is a short sample buffer (32 or 64 bits in size).
type SimpleBuffer struct {
	vec *C.fvec_t
}

// NewSimpleBuffer constructs a new SimpleBuffer.
//
// The caller is responsible for calling Free on the returned
// SimpleBuffer to release memory when done.
//
//     buf := NewSimpleBuffer(bufSize)
//     defer buf.Free()
func NewSimpleBuffer(size uint) *SimpleBuffer {
	return &SimpleBuffer{C.new_fvec(C.uint_t(size))}
}

// NewSimpleBuffer constructs a new SimpleBuffer.
//
// The caller is responsible for calling Free on the returned
// SimpleBuffer to release memory when done.
//
//     buf := NewSimpleBuffer(bufSize)
//     defer buf.Free()
func NewSimpleBufferData(size uint, data []float64) *SimpleBuffer {
	b := C.new_fvec(C.uint_t(size))
	for i := uint(0); i < uint(len(data)); i++ {
		C.fvec_set_sample(b, C.smpl_t(data[i]), C.uint_t(i))
	}
	return &SimpleBuffer{b}
}

// Update the values of the buffer from float64 data
func (b *SimpleBuffer) SetData(data []float64) {
	for i := uint(0); i < uint(len(data)); i++ {
		C.fvec_set_sample(b.vec, C.smpl_t(data[i]), C.uint_t(i))
	}
}

// Update the values of the buffer from float64 data
func (b *SimpleBuffer) SetDataF32(data []float32) {
	for i := C.uint_t(0); i < C.uint_t(len(data)); i++ {
		C.fvec_set_sample(b.vec, C.smpl_t(data[i]), i)
	}
}

// TODO directly operating on the C memory would be fast but also dangerous
// https://copyninja.info/blog/workaround-gotypesystems.html

// Update the values of the buffer using C type casting to avoid conversion overheads.
// This function uses unsafe pointers (careful!)
// The only reason this works is because C.smpl_t is just a float32 under the hood.
// There isn't really any speedup compared to SetDataF32 so no reason to use this.
func (b *SimpleBuffer) SetDataUnsafe(data []float32) {
	cast := *(*[]C.smpl_t)(unsafe.Pointer(&data))
	for i := C.uint_t(0); i < C.uint_t(len(data)); i++ {
		C.fvec_set_sample(b.vec, cast[i], i)
	}
}

// Returns the contents of this buffer as a slice.
// The data is copied so the slices are still valid even
// after the buffer has changed.
func (b *SimpleBuffer) Slice() []float64 {
	sl := make([]float64, b.Size())
	for i := uint(0); i < b.Size(); i++ {
		sl[int(i)] = b.Get(i)
	}
	return sl
}

func (b *SimpleBuffer) Get(i uint) float64 {
	return float64(C.fvec_get_sample(b.vec, C.uint_t(i)))
}

// Size returns the size of this buffer.
func (b *SimpleBuffer) Size() uint {
	if b.vec == nil {
		return 0
	}
	return uint(b.vec.length)
}

// Free frees the memory aubio allocated for this buffer.
func (b *SimpleBuffer) Free() {
	if b.vec == nil {
		return
	}
	C.del_fvec(b.vec)
	b.vec = nil
}

// ComplexBuffer is a wrapper for the aubio cvec_t type.
// It contains complex sample data.
type ComplexBuffer struct {
	data *C.cvec_t
}

// NewComplexBuffer constructs a buffer.
//
// The caller is responsible for calling Free on the returned
// ComplexBuffer to release memory when done.
//
//     buf := NewComplexBuffer(bufSize)
//     defer buf.Free()
func NewComplexBuffer(size uint) *ComplexBuffer {
	return &ComplexBuffer{data: C.new_cvec(C.uint_t(size))}
}

// NewComplexBuffer constructs a buffer with data.
//
func NewComplexBufferData(size uint, data []float64) *ComplexBuffer {
	b := C.new_cvec(C.uint_t(size))
	for i := uint(0); i < uint(len(data)); i++ {
		C.cvec_norm_set_sample(b, C.smpl_t(data[i]), C.uint_t(i))
	}
	return &ComplexBuffer{b}
}

// Free frees the memory aubio has allocated for this buffer.
func (cb *ComplexBuffer) Free() {
	if cb.data != nil {
		C.del_cvec(cb.data)
	}
}

// Size returns the size of this ComplexBuffer.
func (cb *ComplexBuffer) Size() uint {
	if cb.data == nil {
		return 0
	}
	return uint(cb.data.length)
}

// Norm returns the slice of norm data.
// The data is copies so the slice is still
// valid after the buffer has changed.
func (cb *ComplexBuffer) Norm() []float64 {
	sl := make([]float64, cb.Size())
	for i := uint(0); i < cb.Size(); i++ {
		sl[int(i)] = float64(C.cvec_norm_get_sample(cb.data, C.uint_t(i)))
	}
	return sl
}

// Norm returns the slice of phase data.
// The data is copies so the slice is still
// valid after the buffer has changed.
func (cb *ComplexBuffer) Phase() []float64 {
	sl := make([]float64, cb.Size())
	for i := uint(0); i < cb.Size(); i++ {
		sl[int(i)] = float64(C.cvec_phas_get_sample(cb.data, C.uint_t(i)))
	}
	return sl
}

// Buffer for Long sample data (64 bits)
type LongSampleBuffer struct {
	vec *C.lvec_t
}

// NewLBuffer constructs a *LongSampleBuffer.
//
// The caller is responsible for calling Free on the returned
// LongSampleBuffer to release memory when done.
//
//     buf := NewLBuffer(bufSize)
//     defer buf.Free()
func NewLBuffer(size uint) *LongSampleBuffer {
	return newLBufferFromVec(C.new_lvec(C.uint_t(size)))
}

func newLBufferFromVec(v *C.lvec_t) *LongSampleBuffer {
	return &LongSampleBuffer{vec: v}
}

// Free frees the memory allocated by aubio for this buffer.
func (lb *LongSampleBuffer) Free() {
	if lb.vec != nil {
		C.del_lvec(lb.vec)
		lb.vec = nil
	}
}

// Size returns this buffers size.
func (lb *LongSampleBuffer) Size() uint {
	return uint(lb.vec.length)
}

// Returns the contents of this buffer as a slice.
// The data is copied so the slices are still valid even
// after the buffer has changed.
func (lb *LongSampleBuffer) Slice() []float64 {
	sl := make([]float64, lb.Size())
	for i := uint(0); i < lb.Size(); i++ {
		sl[int(i)] = float64(C.lvec_get_sample(lb.vec, C.uint_t(i)))
	}
	return sl
}

type MatrixBuffer struct {
	Length uint
	Height uint
	mat    *C.fmat_t
}

// NewMatBuffer constructs a *MatrixBuffer.
// Height is the number of channels
// Length is the length of a channel
//
// The caller is responsible for calling Free on the returned
// LongSampleBuffer to release memory when done.
//
//     buf := NewLBuffer(bufSize)
//     defer buf.Free()
func NewMatrixBuffer(height, length uint) *MatrixBuffer {
	return NewMatrixBufferFromFmat(C.new_fmat(C.uint_t(height), C.uint_t(length)))
}

func NewMatrixBufferFromFmat(mat *C.fmat_t) *MatrixBuffer {
	return &MatrixBuffer{
		Length: uint(mat.length),
		Height: uint(mat.height),
		mat:    mat,
	}
}

// Free frees the memory allocated by aubio for this buffer.
func (mb *MatrixBuffer) Free() {
	if mb.mat != nil {
		C.del_fmat(mb.mat)
		mb.mat = nil
	}
}

func (mb *MatrixBuffer) Size() uint {
	return mb.Length * mb.Height
}

// Returns the full contents of this matrix buffer as a 2d slice (Length, Height).
// Think of it as a [height]Channel, where Channel is [length]float64
// The data is copied so the slices are still valid even
// after the buffer has changed.
func (mb *MatrixBuffer) GetChannels() [][]float64 {
	sl := make([][]float64, mb.Height)
	for i := range sl {
		sl[i] = make([]float64, mb.Length)
	}
	for i := uint(0); i < mb.Height; i++ {
		for j := uint(0); j < mb.Length; j++ {
			sl[int(i)][int(j)] = float64(C.fmat_get_sample(mb.mat, C.uint_t(i), C.uint_t(j)))
		}
	}
	return sl
}

func (mb *MatrixBuffer) GetChannel(channel uint) []float64 {
	sl := make([]float64, mb.Length)
	for i := uint(0); i < mb.Length; i++ {
		sl[int(i)] = float64(C.fmat_get_sample(mb.mat, C.uint_t(channel), C.uint_t(i)))
	}
	return sl
}
