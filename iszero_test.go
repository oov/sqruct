package sqruct

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

func TestIsZero(t *testing.T) {
	trues := [...]interface{}{
		bool(false),
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		float32(0),
		float64(0),
		"",
		time.Time{},
		(*bool)(nil),
		(*int)(nil),
		(*int8)(nil),
		(*int16)(nil),
		(*int32)(nil),
		(*int64)(nil),
		(*uint)(nil),
		(*uint8)(nil),
		(*uint16)(nil),
		(*uint32)(nil),
		(*uint64)(nil),
		(*float32)(nil),
		(*float64)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		[]byte(nil),
		sql.NullBool{},
		sql.NullInt64{},
		sql.NullFloat64{},
		sql.NullString{},
	}
	for i, v := range trues {
		if !IsZero(v) {
			t.Fatalf("IsZero(trues[%d]) %q: want true got false", i, reflect.TypeOf(v))
		}
	}

	falses := [...]interface{}{
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1),
		float64(1),
		"a",
		time.Now(),
		[]byte(""),
		sql.NullBool{Valid: true},
		sql.NullInt64{Valid: true},
		sql.NullFloat64{Valid: true},
		sql.NullString{Valid: true},
	}
	for i, v := range falses {
		if IsZero(v) {
			t.Fatalf("IsZero(falses[%d]) %q: want false got true", i, reflect.TypeOf(v))
		}
	}
}

func BenchmarkIsZero(b *testing.B) {
	v := sql.NullBool{Valid: true}
	for n := 0; n < b.N; n++ {
		IsZero(v)
	}
}
