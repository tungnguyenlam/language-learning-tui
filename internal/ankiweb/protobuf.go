package ankiweb

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// A minimal protobuf wire-format reader.
//
// AnkiWeb's endpoints answer with protobuf but publish no schema, so the few
// messages this package needs are decoded field by field. Only the wire types
// AnkiWeb actually uses are supported; anything else is skipped so an added
// field in a future response cannot break parsing.
//
// See https://protobuf.dev/programming-guides/encoding/ for the format.

const (
	wireVarint = 0
	wire64Bit  = 1
	wireBytes  = 2
	wire32Bit  = 5
)

var errTruncated = errors.New("truncated protobuf message")

// pbField is one decoded field of a message.
type pbField struct {
	Num   int
	Wire  int
	Value uint64 // for varint / fixed-width fields
	Bytes []byte // for length-delimited fields
}

// pbFields decodes the top-level fields of a protobuf message. Unknown wire
// types end the scan rather than erroring: a partially understood response is
// more useful than none.
func pbFields(data []byte) ([]pbField, error) {
	var fields []pbField
	for len(data) > 0 {
		key, n := binary.Uvarint(data)
		if n <= 0 {
			return fields, errTruncated
		}
		data = data[n:]

		field := pbField{Num: int(key >> 3), Wire: int(key & 7)}
		switch field.Wire {
		case wireVarint:
			v, n := binary.Uvarint(data)
			if n <= 0 {
				return fields, errTruncated
			}
			field.Value, data = v, data[n:]
		case wireBytes:
			length, n := binary.Uvarint(data)
			if n <= 0 {
				return fields, errTruncated
			}
			data = data[n:]
			if uint64(len(data)) < length {
				return fields, errTruncated
			}
			field.Bytes, data = data[:length], data[length:]
		case wire64Bit:
			if len(data) < 8 {
				return fields, errTruncated
			}
			field.Value, data = binary.LittleEndian.Uint64(data[:8]), data[8:]
		case wire32Bit:
			if len(data) < 4 {
				return fields, errTruncated
			}
			field.Value, data = uint64(binary.LittleEndian.Uint32(data[:4])), data[4:]
		default:
			return fields, fmt.Errorf("unsupported protobuf wire type %d", field.Wire)
		}
		fields = append(fields, field)
	}
	return fields, nil
}

// pbSubMessage returns the bytes of the first length-delimited field with the
// given number.
func pbSubMessage(fields []pbField, num int) ([]byte, bool) {
	for _, f := range fields {
		if f.Num == num && f.Wire == wireBytes {
			return f.Bytes, true
		}
	}
	return nil, false
}

// pbString returns the first length-delimited field with the given number as a
// string.
func pbString(fields []pbField, num int) string {
	if b, ok := pbSubMessage(fields, num); ok {
		return string(b)
	}
	return ""
}

// pbUint returns the first varint field with the given number.
func pbUint(fields []pbField, num int) uint64 {
	for _, f := range fields {
		if f.Num == num && f.Wire == wireVarint {
			return f.Value
		}
	}
	return 0
}
