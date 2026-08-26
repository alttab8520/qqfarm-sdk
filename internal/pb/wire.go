package pb

import "google.golang.org/protobuf/encoding/protowire"

type Encoder struct {
	buf []byte
}

func NewEncoder() *Encoder { return &Encoder{} }

func (e *Encoder) Bytes() []byte { return e.buf }

func (e *Encoder) Varint(num protowire.Number, v uint64) {
	e.buf = protowire.AppendTag(e.buf, num, protowire.VarintType)
	e.buf = protowire.AppendVarint(e.buf, v)
}

func (e *Encoder) Int(num protowire.Number, v int64) {
	if v == 0 {
		return
	}
	e.Varint(num, uint64(v))
}

func (e *Encoder) Bool(num protowire.Number, v bool) {
	if !v {
		return
	}
	e.Varint(num, 1)
}

func (e *Encoder) String(num protowire.Number, s string) {
	if s == "" {
		return
	}
	e.BytesField(num, []byte(s))
}

func (e *Encoder) BytesField(num protowire.Number, b []byte) {
	if len(b) == 0 {
		return
	}
	e.buf = protowire.AppendTag(e.buf, num, protowire.BytesType)
	e.buf = protowire.AppendBytes(e.buf, b)
}

func (e *Encoder) Message(num protowire.Number, inner []byte) {
	e.BytesField(num, inner)
}

func (e *Encoder) PackedVarints(num protowire.Number, vs []int64) {
	if len(vs) == 0 {
		return
	}
	var packed []byte
	for _, v := range vs {
		packed = protowire.AppendVarint(packed, uint64(v))
	}
	e.BytesField(num, packed)
}

func (e *Encoder) RepeatedVarint(num protowire.Number, vs []int64) {
	for _, v := range vs {
		e.Varint(num, uint64(v))
	}
}

type Field struct {
	Num    protowire.Number
	Varint uint64
	Bytes  []byte
	Kind   protowire.Type
}

func Walk(data []byte) ([]Field, error) {
	var out []Field
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return out, protowire.ParseError(n)
		}
		data = data[n:]
		f := Field{Num: num, Kind: typ}
		switch typ {
		case protowire.VarintType:
			v, m := protowire.ConsumeVarint(data)
			if m < 0 {
				return out, protowire.ParseError(m)
			}
			f.Varint = v
			data = data[m:]
		case protowire.BytesType:
			b, m := protowire.ConsumeBytes(data)
			if m < 0 {
				return out, protowire.ParseError(m)
			}
			f.Bytes = b
			data = data[m:]
		case protowire.Fixed32Type:
			_, m := protowire.ConsumeFixed32(data)
			if m < 0 {
				return out, protowire.ParseError(m)
			}
			data = data[m:]
		case protowire.Fixed64Type:
			_, m := protowire.ConsumeFixed64(data)
			if m < 0 {
				return out, protowire.ParseError(m)
			}
			data = data[m:]
		default:
			return out, protowire.ParseError(int(typ))
		}
		out = append(out, f)
	}
	return out, nil
}

func FieldMap(data []byte) (map[protowire.Number]Field, error) {
	fields, err := Walk(data)
	if err != nil {
		return nil, err
	}
	out := make(map[protowire.Number]Field, len(fields))
	for _, f := range fields {
		out[f.Num] = f
	}
	return out, nil
}

func IntField(m map[protowire.Number]Field, num protowire.Number) int64 {
	f, ok := m[num]
	if !ok {
		return 0
	}
	return int64(f.Varint)
}

func StringField(m map[protowire.Number]Field, num protowire.Number) string {
	f, ok := m[num]
	if !ok {
		return ""
	}
	return string(f.Bytes)
}

func BoolField(m map[protowire.Number]Field, num protowire.Number) bool {
	return IntField(m, num) != 0
}
