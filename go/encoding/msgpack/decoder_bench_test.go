package msgpack

// $ go test *.go -test.bench Decode -test.benchmem -test.cpu=1,4,8 -test.benchtime=10s
// Benchmark__Msgpack__Decode_Small      5000000     4313 ns/op     963 B/op   26 allocs/op
// Benchmark__Msgpack__Decode_Small-4   10000000     1911 ns/op     960 B/op   26 allocs/op
// Benchmark__Msgpack__Decode_Small-8   10000000     1656 ns/op     958 B/op   26 allocs/op
// Benchmark__Msgpack__Decode_Large       100000   251495 ns/op   29620 B/op   27 allocs/op
// Benchmark__Msgpack__Decode_Large-4     200000    87714 ns/op   29617 B/op   27 allocs/op
// Benchmark__Msgpack__Decode_Large-8     200000    84622 ns/op   29617 B/op   27 allocs/op
// Benchmark__Pool_____Decode_Small      5000000     4189 ns/op     753 B/op   25 allocs/op
// Benchmark__Pool_____Decode_Small-4   10000000     1656 ns/op     750 B/op   25 allocs/op
// Benchmark__Pool_____Decode_Small-8   10000000     1581 ns/op     747 B/op   25 allocs/op
// Benchmark__Pool_____Decode_Large       100000   254045 ns/op   29409 B/op   26 allocs/op
// Benchmark__Pool_____Decode_Large-4     200000    88093 ns/op   29409 B/op   26 allocs/op
// Benchmark__Pool_____Decode_Large-8     200000    83745 ns/op   29409 B/op   26 allocs/op

import (
	"runtime"
	"sync"
	"testing"

	"github.com/ugorji/go/codec"
)

var bSmall, bLarge []byte

func init() {
	bSmall, _ = Marshal(small)
	bLarge, _ = Marshal(large)
}

type benchDecFn func(b *testing.B, buf []byte, n int)

func benchDec(b *testing.B, buf []byte, fn benchDecFn) {
	runtime.GC()
	b.ResetTimer()
	var (
		w sync.WaitGroup
		m = runtime.GOMAXPROCS(0)
		n = b.N / m
	)
	w.Add(m)
	for i := 0; i < m; i++ {
		go func() {
			fn(b, buf, n)
			w.Done()
		}()
	}
	w.Wait()
}

func benchDecMsgpack(b *testing.B, buf []byte, n int) {
	h := &codec.MsgpackHandle{}
	for i := 0; i < n; i++ {
		var v map[string]interface{}
		if err := codec.NewDecoderBytes(buf, h).Decode(&v); err != nil {
			b.Log("Error decoding:", err)
			b.FailNow()
		}
	}
}

func benchDecPool(b *testing.B, buf []byte, n int) {
	for i := 0; i < n; i++ {
		var v map[string]interface{}
		if err := Unmarshal(buf, &v); err != nil {
			b.Log("Error decoding:", err)
			b.FailNow()
		}
	}
}

func Benchmark__Msgpack__Decode_Small(b *testing.B) {
	benchDec(b, bSmall, benchDecMsgpack)
}

func Benchmark__Msgpack__Decode_Large(b *testing.B) {
	benchDec(b, bLarge, benchDecMsgpack)
}

func Benchmark__Pool_____Decode_Small(b *testing.B) {
	benchDec(b, bSmall, benchDecPool)
}

func Benchmark__Pool_____Decode_Large(b *testing.B) {
	benchDec(b, bLarge, benchDecPool)
}
