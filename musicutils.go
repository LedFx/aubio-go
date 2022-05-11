package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

// Compute sound pressure level (SPL) in dB
func db_spl(buf *SimpleBuffer) float64 {
	return float64(C.aubio_db_spl(buf.vec))
}

// Get level in dB SPL if level >= threshold, otherwise 1.
func LevelDetection(buf *SimpleBuffer, threshold float64) float64 {
	return float64(C.aubio_level_detection(buf.vec, C.smpl_t(threshold)))
}

// Compute sound level on a linear scale
func LevelLin(buf *SimpleBuffer) float64 {
	return float64(C.aubio_level_lin(buf.vec))
}

// Convert frequency (Hz) to mel.
// Converts a scalar from the frequency domain to the mel scale using Slaney Auditory Toolbox's implementation
func HzToMel(freq float64) float64 {
	return float64(C.aubio_hztomel(C.smpl_t(freq)))
}

// Convert frequency (Hz) to mel using HTK scaling.
// Converts a scalar from the frequency domain to the mel scale, using the equation defined by O'Shaughnessy, as implemented in the HTK speech recognition toolkit
func HzToMelHTK(freq float64) float64 {
	return float64(C.aubio_hztomel_htk(C.smpl_t(freq)))
}

// Convert mel to frequency (Hz).
// Converts a scalar from the mel scale to the frequency domain using Slaney Auditory Toolbox's implementation
func MelToHz(mel float64) float64 {
	return float64(C.aubio_meltohz(C.smpl_t(mel)))
}

// Convert mel to frequency (Hz) using HTK scaling.
// Converts a scalar from the mel scale to the frequency domain, using the equation defined by O'Shaughnessy, as implemented in the HTK speech recognition toolkit
func MelToHzHTK(mel float64) float64 {
	return float64(C.aubio_meltohz_htk(C.smpl_t(mel)))
}

// Check if buffer level in dB SPL is under a given threshold
// True if level is under threshold, false otherwise
func SilenceDetection(buf *SimpleBuffer, threshold float64) bool {
	return C.aubio_silence_detection(buf.vec, C.smpl_t(threshold)) == 0
}

//Compute the principal argument.
//This function maps the input phase to its corresponding value wrapped in the range [−π,π].
func Unwrap2pi(phase float64) float64 {
	return float64(C.aubio_unwrap2pi(C.smpl_t(phase)))
}

// The zero-crossing rate is the number of times a signal changes sign, divided by the length of this signal.
func ZeroCrossingRate(buf *SimpleBuffer) float64 {
	return float64(C.aubio_zero_crossing_rate(buf.vec))
}
