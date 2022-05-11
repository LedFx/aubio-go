package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

// Inplace compute the exp(x) of each vector elements
func (buf *SimpleBuffer) Exp() {
	C.fvec_exp(buf.vec)
}

// Inplace compute the cos(x) of each vector elements
func (buf *SimpleBuffer) Cos() {
	C.fvec_cos(buf.vec)
}

// Inplace compute the sin(x) of each vector elements
func (buf *SimpleBuffer) Sin() {
	C.fvec_sin(buf.vec)
}

// Inplace compute the abs(x) of each vector elements
func (buf *SimpleBuffer) Abs() {
	C.fvec_abs(buf.vec)
}

// Inplace compute the sqrt(x) of each vector elements
func (buf *SimpleBuffer) Sqrt() {
	C.fvec_sqrt(buf.vec)
}

// Inplace compute the log10(x) of each vector elements
func (buf *SimpleBuffer) Log10() {
	C.fvec_log10(buf.vec)
}

// Inplace compute the log(x) (natural log, ln) of each vector elements
func (buf *SimpleBuffer) Log() {
	C.fvec_log(buf.vec)
}

// Inplace compute the floor(x) of each vector elements
func (buf *SimpleBuffer) Floor() {
	C.fvec_floor(buf.vec)
}

// Inplace compute the ceil(x) of each vector elements
func (buf *SimpleBuffer) Ceil() {
	C.fvec_ceil(buf.vec)
}

// Inplace compute the round(x) of each vector elements
func (buf *SimpleBuffer) Round() {
	C.fvec_round(buf.vec)
}

// Inplace raise each vector elements to the power pow
func (buf *SimpleBuffer) Pow(pow float64) {
	C.fvec_pow(buf.vec, C.smpl_t(pow))
}

// Inplace clamp the values of a buffer within the range [-abs(max), abs(max)]
func (buf *SimpleBuffer) Clamp(absmax float64) {
	C.fvec_clamp(buf.vec, C.smpl_t(absmax))
}
