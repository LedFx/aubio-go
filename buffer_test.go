package aubio

import "testing"

func TestBufferWriteFast(t *testing.T) {
	b := NewSimpleBuffer(100)
	data := make([]float64, 100)
	b.SetDataFast(data)
}
