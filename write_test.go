package etf

import (
	"bytes"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestWriteAtom(t *testing.T) {
	c := new(Context)
	test := func(in Atom, shouldFail bool) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeAtom(in); err != nil {
			if !shouldFail {
				t.Error(in, err)
			}
		} else if shouldFail {
			t.Errorf("err == nil (%v)", in)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if v != in {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test(Atom(""), false)
	test(Atom(bytes.Repeat([]byte{'a'}, math.MaxUint8)), false)
	test(Atom(bytes.Repeat([]byte{'a'}, math.MaxUint8+1)), false)
	test(Atom(bytes.Repeat([]byte{'a'}, math.MaxUint16)), false)
	test(Atom(bytes.Repeat([]byte{'a'}, math.MaxUint16+1)), true)
}

func TestWriteBinary(t *testing.T) {
	c := new(Context)
	test := func(in []byte) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeBinary(in); err != nil {
			t.Error(in, err)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if bytes.Compare(v.([]byte), in) != 0 {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test([]byte{})
	test(bytes.Repeat([]byte{231}, 65535))
	test(bytes.Repeat([]byte{123}, 65536))
}

func TestWriteBool(t *testing.T) {
	c := new(Context)
	test := func(in bool) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeBool(in); err != nil {
			t.Error(in, err)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if v != in {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test(true)
	test(false)
}

func TestWriteFloat(t *testing.T) {
	c := new(Context)
	test := func(in float64) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeFloat(in); err != nil {
			t.Error(in, err)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if v != in {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test(0.0)
	test(-12345.6789)
	test(math.SmallestNonzeroFloat64)
	test(math.MaxFloat64)
}

func TestWriteInt(t *testing.T) {
	c := new(Context)
	test := func(in int64) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeInt(in); err != nil {
			t.Error(in, err)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if reflect.ValueOf(v).Int() != in {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test(0)
	test(-1)
	test(math.MaxInt8)
	test(math.MaxInt8 + 1)
	test(math.MaxInt32)
	test(math.MaxInt32 + 1)
	test(math.MinInt32)
	test(math.MinInt32 - 1)
	test(math.MinInt64)
	test(math.MaxInt64)
}

func TestWriteUint(t *testing.T) {
	c := new(Context)
	test := func(in uint64) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeUint(in); err != nil {
			t.Error(in, err)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else {
			var uv uint64
			switch v := v.(type) {
			case int:
				uv = uint64(v)
			case int64:
				uv = uint64(v)
			case *big.Int:
				uv = v.Uint64()
			}
			if uv != in {
				t.Errorf("expected %v, got %v", in, v)
			}
		}
	}

	test(0)
	test(math.MaxUint8)
	test(math.MaxUint8 + 1)
	test(math.MaxUint32)
	test(math.MaxUint32 + 1)
	test(math.MaxUint64)
}

func TestWritePid(t *testing.T) {
	c := new(Context)
	test := func(in Pid) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writePid(in); err != nil {
			t.Error(in, err)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if v != in {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test(Pid{Atom("omg@lol"), 38, 0, 3})
	test(Pid{Atom("self@localhost"), 32, 1, 9})
}

func TestWriteString(t *testing.T) {
	c := new(Context)
	test := func(in string, shouldFail bool) {
		w := new(bytes.Buffer)
		e := c.Encoder(w)
		if err := e.writeString(in); err != nil {
			if !shouldFail {
				t.Error(in, err)
			}
		} else if shouldFail {
			t.Errorf("err == nil (%v)", in)
		} else if v, err := c.Decoder(w).Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if v != in {
			t.Errorf("expected %v, got %v", in, v)
		}
	}

	test(string(bytes.Repeat([]byte{'a'}, math.MaxUint16)), false)
	test("", false)
	test(string(bytes.Repeat([]byte{'a'}, math.MaxUint16+1)), true)
}

func TestWriteTerm(t *testing.T) {
	c := new(Context)
	type s1 struct {
		L []any
		F float64
	}
	type s2 struct {
		Atom
		S  string
		I  int
		S1 s1
		B  byte
	}
	in := s2{
		Atom("lol"),
		"omg",
		13666,
		s1{
			[]any{
				256,
				"1",
				13.0,
			},
			13.13,
		},
		1,
	}

	w := new(bytes.Buffer)
	e := c.Encoder(w)
	if err := e.Encode(in); err != nil {
		t.Error(in, err)
	} else {
		d := c.Decoder(w)
		if term, err := d.Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", in, l)
		} else if err := e.Encode(term); err != nil {
			t.Error(term, err)
		} else if term, err := d.Decode(); err != nil {
			t.Error(in, err)
		} else if l := w.Len(); l != 0 {
			t.Errorf("%v: buffer len %d", term, l)
		}
	}
}
