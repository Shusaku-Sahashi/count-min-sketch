package src

// ref: https://github.com/dgraph-io/ristretto/blob/29b4dd7a077785696ba5422081b3c301d4c1e5f1/sketch.go
const (
	depth = 4
	shift = 1 << 5
)

type Sketch struct {
	data [depth]rowElm
	mask uint64
}

func NewSketch(counterSize int64) *Sketch {
	if counterSize == 0 {
		panic("msSketch: bad counterSize")
	}
	counterSize = next2Power(counterSize)
	sketch := &Sketch{mask: uint64(counterSize - 1)}
	for i := 0; i < depth; i++ {
		sketch.data[i] = newRowElem(counterSize)
	}
	return sketch
}

func nHash(h, l, i, m uint64) uint64{
	return (h + i * l) % m
}

func (s Sketch) Increment(hash uint64) {
	h := hash >> shift
	l := hash << shift >> shift
	for i, d := range s.data {
		d.increment(nHash(h, l, uint64(i), s.mask))
	}
}

func (s Sketch) Estimate(hash uint64) int64{
	min := byte(255)
	h := hash >> shift
	l := hash << shift >> shift
	for i, d := range s.data {
		v := d.get(nHash(h, l, uint64(i), s.mask))
		if min > v {
			min = v
		}
	}

	return int64(min)
}

/*
	1111       |1111      8bit
	even area  |odd area
 */
type rowElm []byte

func newRowElem(num int64) rowElm {
	return make([]byte, num/2)
}

func (r rowElm) get(n uint64) byte{
	return r[n/2]>>((n&1)*4)
}

func (r rowElm) increment(n uint64) {
	idx := n / 2

	s := (n & 1) * 4 // odd => +1, even => +10000
	v := (r[idx] >> s) & 0x0f

	// upper limit of count is 15
	if v < 15 {
		r[idx] += 1 << s
	}
}

// generate next 2 pow number.
func next2Power(x int64) int64 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}

