package etf

import (
	"bytes"
	"io"
	"math/rand/v2"
	"testing"
)

func BenchmarkWriteAtom(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 64
	length := 64
	atoms := make([]Atom, max)

	for i := 0; i < max; i++ {
		atoms[i] = Atom(bytes.Repeat([]byte{byte('A' + i)}, length))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := atoms[i%max]
		if err := e.writeAtom(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

//func BenchmarkWriteBigInt(b *testing.B) {
//	b.StopTimer()
//	c := new(Context)
//
//	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
//	uint64Max := new(big.Int).SetUint64(math.MaxUint64)
//	top := new(big.Int).Mul(uint64Max, uint64Max)
//	max := 512
//	bigints := make([]*big.Int, max)
//
//	for i := 0; i < max; i++ {
//		a := new(big.Int).Rand(rand, top)
//		b := new(big.Int).Rand(rand, top)
//		bigints[i] = new(big.Int).Sub(a, b)
//	}
//
//	b.StartTimer()
//
//	for i := 0; i < b.N; i++ {
//		in := bigints[i%max]
//		if err := c.writeBigInt(Discard, in); err != nil {
//			b.Fatal(in, err)
//		}
//	}
//}

func BenchmarkWriteBinary(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 64
	length := 64
	binaries := make([][]byte, max)

	for i := 0; i < max; i++ {
		s := bytes.Repeat([]byte{'a'}, length)
		binaries[i] = bytes.Map(
			func(rune) rune { return rune(byte(rand.Int())) },
			s,
		)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := binaries[i%max]
		if err := e.writeBinary(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

func BenchmarkWriteBool(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 64
	bools := make([]bool, max)

	for i := 0; i < max; i++ {
		bools[i] = (rand.N(2) == 1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := bools[i%max]
		if err := e.writeBool(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

func BenchmarkWriteFloat(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 512
	floats := make([]float64, max)

	for i := 0; i < max; i++ {
		floats[i] = rand.ExpFloat64() - rand.ExpFloat64()
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := floats[i%max]
		if err := e.writeFloat(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

func BenchmarkWriteInt(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 512
	ints := make([]int64, max)

	for i := 0; i < max; i++ {
		ints[i] = int64(rand.Int32() - rand.Int32())
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := ints[i%max]
		if err := e.writeInt(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

func BenchmarkWriteUint(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 512
	ints := make([]uint64, max)

	for i := 0; i < max; i++ {
		ints[i] = uint64(rand.Int32())
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := ints[i%max]
		if err := e.writeUint(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

func BenchmarkWritePid(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 64
	length := 16
	pids := make([]Pid, max)

	for i := 0; i < max; i++ {
		s := bytes.Repeat([]byte{'a'}, length)
		b := bytes.Map(randRune, s)
		b[6] = '@'
		pids[i] = Pid{
			Atom(b),
			rand.N[uint32](65536),
			rand.N[uint32](256),
			rand.N[byte](16),
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := pids[i%max]
		if err := e.writePid(in); err != nil {
			b.Fatal(in, err)
		}
	}
}

func BenchmarkWriteString(b *testing.B) {
	b.StopTimer()
	c := new(Context)
	e := c.Encoder(io.Discard)

	max := 64
	length := 64
	strings := make([]string, max)

	for i := 0; i < max; i++ {
		s := bytes.Repeat([]byte{'a'}, length)
		strings[i] = string(bytes.Map(randRune, s))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		in := strings[i%max]
		if err := e.writeString(in); err != nil {
			b.Fatal(in, err)
		}
	}
}
