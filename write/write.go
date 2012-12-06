package write

import (
	"fmt"
	. "github.com/goerlang/etf/types"
	"io"
	"math"
	"math/big"
)

func Atom(w io.Writer, atom ErlAtom) (err error) {
	switch size := len(atom); {
	case size <= 0xff:
		// $sL…
		if _, err = w.Write([]byte{ErlTypeSmallAtom, byte(size)}); err == nil {
			_, err = w.Write([]byte(atom))
		}

	case size <= 0xffff:
		// $dLL…
		_, err = w.Write([]byte{byte(ErlTypeAtom), byte(size >> 8), byte(size)})
		if err == nil {
			_, err = w.Write([]byte(atom))
		}

	default:
		err = fmt.Errorf("atom is too big (%d bytes)", size)
	}

	return
}

func BigInt(w io.Writer, x *big.Int) (err error) {
	sign := 0
	if x.Sign() < 0 {
		sign = 1
	}

	bytes := reverse(new(big.Int).Abs(x).Bytes())

	switch size := len(bytes); {
	case size <= 0xff:
		// $nAS…
		_, err = w.Write([]byte{ErlTypeSmallBig, byte(size), byte(sign)})
		if err == nil {
			_, err = w.Write(bytes)
		}

	case int(uint32(size)) == size:
		// $oAAAAS…
		data := []byte{
			ErlTypeLargeBig,
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
			byte(sign),
		}
		if _, err = w.Write(data); err == nil {
			_, err = w.Write(bytes)
		}

	default:
		err = fmt.Errorf("bad big int size (%d)", size)
	}

	return
}

func Binary(w io.Writer, bytes []byte) (err error) {
	switch size := len(bytes); {
	case int(uint32(size)) == size:
		// $mLLLL…
		data := []byte{
			ErlTypeBinary,
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		}
		if _, err = w.Write(data); err == nil {
			_, err = w.Write(bytes)
		}

	default:
		err = fmt.Errorf("bad binary size (%d)", size)
	}

	return
}

func Bool(w io.Writer, b bool) (err error) {
	switch b {
	case true:
		err = Atom(w, ErlAtom("true"))

	case false:
		err = Atom(w, ErlAtom("false"))
	}

	return
}

func Float64(w io.Writer, f float64) (err error) {
	if _, err = w.Write([]byte{ErlTypeNewFloat}); err == nil {
		fb := math.Float64bits(f)
		_, err = w.Write([]byte{
			byte(fb >> 56), byte(fb >> 48), byte(fb >> 40), byte(fb >> 32),
			byte(fb >> 24), byte(fb >> 16), byte(fb >> 8), byte(fb),
		})
	}
	return
}

func Int64(w io.Writer, x int64) (err error) {
	switch {
	case x >= 0 && x <= 0xff:
		// $aI
		_, err = w.Write([]byte{ErlTypeSmallInteger, byte(x)})

	case int64(int32(x)) == x:
		// $bIIII
		_, err = w.Write([]byte{
			ErlTypeInteger,
			byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x),
		})

	default:
		err = BigInt(w, big.NewInt(x))
	}

	return
}

func String(w io.Writer, s string) (err error) {
	switch size := len(s); {
	case size <= 0xffff:
		// $kLL…
		_, err = w.Write([]byte{ErlTypeString, byte(size >> 8), byte(size)})
		if err == nil {
			_, err = w.Write([]byte(s))
		}

	default:
		err = fmt.Errorf("string is too big (%d bytes)", size)
	}

	return
}

func reverse(b []byte) []byte {
	size := len(b)
	r := make([]byte, size)

	for i := 0; i < size; i++ {
		r[i] = b[size-i-1]
	}

	return r
}
