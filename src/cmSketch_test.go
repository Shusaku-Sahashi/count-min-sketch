package src

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSketch(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()

	s := NewSketch(5)
	require.Equal(t, uint64(7), s.mask)
	NewSketch(0)
}

func TestSketchEstimate(t *testing.T) {
	s := NewSketch(16)
	s.Increment(1)
	s.Increment(1)
	require.Equal(t, int64(2), s.Estimate(1))
	require.Equal(t, int64(0), s.Estimate(0))
}
func TestNext2Power(t *testing.T) {
	sz := 12 << 30
	szf := float64(sz) * 0.01
	val := int64(szf)
	t.Logf("szf = %.2f val = %d\n", szf, val)
	pow := next2Power(val)
	t.Logf("pow = %d. mult 4 = %d\n", pow, pow*4)
}

func BenchmarkSketchIncrement(b *testing.B) {
	s := NewSketch(16)
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		s.Increment(1)
	}
}

func BenchmarkSketchEstimate(b *testing.B) {
	s := NewSketch(16)
	s.Increment(1)
	b.SetBytes(1)
	for n := 0; n < b.N; n++ {
		s.Estimate(1)
	}
}

