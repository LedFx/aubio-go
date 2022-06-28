package aubio

import (
	"fmt"
	"testing"
)

func TestSimpleBufferSetFast(t *testing.T) {
	b := NewSimpleBuffer(100)
	data := make([]float64, 100)
	b.SetDataFast(data)
}

func BenchmarkSimpleBuffer(t *testing.B) {
	b := NewSimpleBuffer(100)
	lens := []int{50, 100, 500, 1000, 5000}
	for _, l := range lens {
		data := make([]float64, l)
		t.Run(fmt.Sprintf("%v set data slow", l), func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				b.SetData(data)
			}
		})
		t.Run(fmt.Sprintf("%v set data fast", l), func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				b.SetDataFast(data)
			}
		})
	}
}
