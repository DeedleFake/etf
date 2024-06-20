package etf

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
)

type Encoder struct {
	c *Context
	w io.Writer
}

func (e *Encoder) Encode(term interface{}) (err error) {
	_, err = e.w.Write([]byte{EtVersion})
	if err != nil {
		return err
	}

	return e.EncodeTerm(term)
}

func (e *Encoder) EncodeTerm(term any) (err error) {
	switch v := term.(type) {
	case bool:
		err = e.writeBool(v)
	case int8, int16, int32, int64, int:
		err = e.writeInt(reflect.ValueOf(term).Int())
	case uint8, uint16, uint32, uint64, uintptr, uint:
		err = e.writeUint(reflect.ValueOf(term).Uint())
	case *big.Int:
		err = e.writeBigInt(v)
	case string:
		err = e.writeString(v)
	case []byte:
		err = e.writeBinary(v)
	case float64:
		err = e.writeFloat(v)
	case float32:
		err = e.writeFloat(float64(v))
	case Atom:
		err = e.writeAtom(v)
	case Pid:
		err = e.writePid(v)
	case Tuple:
		err = e.writeTuple(v)
	case Ref:
		err = e.writeRef(v)
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Struct:
			err = e.writeRecord(term)
		case reflect.Array, reflect.Slice:
			err = e.writeList(term)
		case reflect.Ptr:
			err = e.EncodeTerm(rv.Elem())
		//case reflect.Map // FIXME
		default:
			err = &ErrUnknownType{rv.Type()}
		}
	}

	return
}

func (e *Encoder) writeAtom(atom Atom) (err error) {
	switch size := len(atom); {
	case size <= math.MaxUint8:
		// $sL…
		if _, err = e.w.Write([]byte{ettSmallAtom, byte(size)}); err == nil {
			_, err = io.WriteString(e.w, string(atom))
		}

	case size <= math.MaxUint16:
		// $dLL…
		_, err = e.w.Write([]byte{ettAtom, byte(size >> 8), byte(size)})
		if err == nil {
			_, err = io.WriteString(e.w, string(atom))
		}

	default:
		err = fmt.Errorf("atom is too big (%d bytes)", size)
	}

	return
}

func (e *Encoder) writeBigInt(x *big.Int) (err error) {
	sign := 0
	if x.Sign() < 0 {
		sign = 1
	}

	bytes := reverse(new(big.Int).Abs(x).Bytes())

	switch size := int64(len(bytes)); {
	case size <= math.MaxUint8:
		// $nAS…
		_, err = e.w.Write([]byte{ettSmallBig, byte(size), byte(sign)})

	case size <= math.MaxUint32:
		// $oAAAAS…
		_, err = e.w.Write([]byte{
			ettLargeBig,
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
			byte(sign),
		})

	default:
		err = fmt.Errorf("bad big int size (%d)", size)
	}

	if err == nil {
		_, err = e.w.Write(bytes)
	}

	return
}

func (e *Encoder) writeBinary(bytes []byte) (err error) {
	switch size := int64(len(bytes)); {
	case size <= math.MaxUint32:
		// $mLLLL…
		data := []byte{
			ettBinary,
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		}
		if _, err = e.w.Write(data); err == nil {
			_, err = e.w.Write(bytes)
		}

	default:
		err = fmt.Errorf("bad binary size (%d)", size)
	}

	return
}

func (e *Encoder) writeBool(b bool) (err error) {
	// $sL…
	if b {
		_, err = e.w.Write([]byte{ettSmallAtom, 4, 't', 'r', 'u', 'e'})
	} else {
		_, err = e.w.Write([]byte{ettSmallAtom, 5, 'f', 'a', 'l', 's', 'e'})
	}

	return
}

func (e *Encoder) writeFloat(f float64) (err error) {
	if _, err = e.w.Write([]byte{ettNewFloat}); err == nil {
		fb := math.Float64bits(f)
		_, err = e.w.Write([]byte{
			byte(fb >> 56), byte(fb >> 48), byte(fb >> 40), byte(fb >> 32),
			byte(fb >> 24), byte(fb >> 16), byte(fb >> 8), byte(fb),
		})
	}
	return
}

func (e *Encoder) writeInt(x int64) (err error) {
	switch {
	case x >= 0 && x <= math.MaxUint8:
		// $aI
		_, err = e.w.Write([]byte{ettSmallInteger, byte(x)})

	case x >= math.MinInt32 && x <= math.MaxInt32:
		// $bIIII
		x := int32(x)
		_, err = e.w.Write([]byte{
			ettInteger,
			byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x),
		})

	default:
		err = e.writeBigInt(big.NewInt(x))
	}

	return
}

func (e *Encoder) writeUint(x uint64) (err error) {
	switch {
	case x <= math.MaxUint8:
		// $aI
		_, err = e.w.Write([]byte{ettSmallInteger, byte(x)})

	case x <= math.MaxInt32:
		// $bIIII
		_, err = e.w.Write([]byte{
			ettInteger,
			byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x),
		})

	default:
		err = e.writeBigInt(new(big.Int).SetUint64(x))
	}

	return
}

func (e *Encoder) writePid(p Pid) (err error) {
	if _, err = e.w.Write([]byte{ettPid}); err != nil {
		return
	} else if err = e.writeAtom(p.Node); err != nil {
		return
	}

	_, err = e.w.Write([]byte{
		0, 0, byte(p.Id >> 8), byte(p.Id),
		byte(p.Serial >> 24),
		byte(p.Serial >> 16),
		byte(p.Serial >> 8),
		byte(p.Serial),
		p.Creation,
	})

	return
}

func (e *Encoder) writeString(s string) (err error) {
	switch size := len(s); {
	case size <= math.MaxUint16:
		// $kLL…
		_, err = e.w.Write([]byte{ettString, byte(size >> 8), byte(size)})
		if err == nil {
			_, err = e.w.Write([]byte(s))
		}

	default:
		err = fmt.Errorf("string is too big (%d bytes)", size)
	}

	return
}

func (e *Encoder) writeList(l interface{}) (err error) {
	rv := reflect.ValueOf(l)
	n := rv.Len()
	_, err = e.w.Write([]byte{
		ettList,
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})

	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		v := rv.Index(i).Interface()
		if err = e.EncodeTerm(v); err != nil {
			return
		}
	}

	_, err = e.w.Write([]byte{ettNil})

	return
}

func (e *Encoder) writeRecord(r interface{}) (err error) {
	rv := reflect.ValueOf(r)
	rt := rv.Type()
	fields := make([]reflect.StructField, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous || !field.IsExported() {
			continue
		}
		fields = append(fields, field)
	}

	if len(fields) <= math.MaxUint8 {
		_, err = e.w.Write([]byte{ettSmallTuple, byte(len(fields))})
	} else {
		_, err = e.w.Write([]byte{
			ettLargeTuple,
			byte(len(fields) >> 24),
			byte(len(fields) >> 16),
			byte(len(fields) >> 8),
			byte(len(fields)),
		})
	}
	if err != nil {
		return err
	}

	for _, field := range fields {
		f := rv.FieldByIndex(field.Index)
		if err = e.EncodeTerm(f.Interface()); err != nil {
			return err
		}
	}

	return err
}

func (e *Encoder) writeRef(ref Ref) (err error) {
	n := len(ref.Id)
	_, err = e.w.Write([]byte{ettNewRef, byte(n >> 8), byte(n)})
	if err != nil {
		return
	}
	if err = e.writeAtom(ref.Node); err != nil {
		return
	}
	if _, err = e.w.Write([]byte{ref.Creation}); err != nil {
		return
	}
	for _, v := range ref.Id {
		b := []byte{
			byte(v >> 24),
			byte(v >> 16),
			byte(v >> 8),
			byte(v),
		}
		if _, err = e.w.Write(b); err != nil {
			return
		}
	}

	return
}

func (e *Encoder) writeTuple(tuple Tuple) (err error) {
	n := len(tuple)
	if n <= math.MaxUint8 {
		_, err = e.w.Write([]byte{ettSmallTuple, byte(n)})
	} else {
		_, err = e.w.Write([]byte{
			ettLargeTuple,
			byte(n >> 24),
			byte(n >> 16),
			byte(n >> 8),
			byte(n),
		})
	}

	if err != nil {
		return
	}

	for _, v := range tuple {
		if err = e.EncodeTerm(v); err != nil {
			return
		}
	}

	return
}

// ErrUnknownType is returned by an attempt to write a type that isn't
// supported.
type ErrUnknownType struct {
	t reflect.Type
}

func (e *ErrUnknownType) Error() string {
	return fmt.Sprintf("write: can't encode type \"%s\"", e.t.Name())
}
