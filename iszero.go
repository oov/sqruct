package sqruct

import (
	"database/sql"
	"reflect"
	"time"
)

// IsZero reports whether x is zero value. This supports following types.
//	bool, *bool
//	int, *int
//	int8, *int8
//	int16, *int16
//	int32, *int32
//	int64, *int64
//	uint, *uint
//	uint8, *uint8
//	uint16, *uint16
//	uint32, *uint32
//	uint64, *uint64
//	float32, *float32
//	float64, *float64
//	string, *string
//	[]byte
//	time.Time, *time.Time
//	sql.NullBool
//	sql.NullInt64
//	sql.NullFloat64
//	sql.NullString
func IsZero(x interface{}) bool {
	switch v := x.(type) {
	case bool:
		return v == false
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case string:
		return len(v) == 0
	case time.Time:
		return v.IsZero()
	case sql.NullBool:
		return !v.Valid
	case sql.NullInt64:
		return !v.Valid
	case sql.NullFloat64:
		return !v.Valid
	case sql.NullString:
		return !v.Valid
	case *bool:
		return v == nil
	case *int:
		return v == nil
	case *int8:
		return v == nil
	case *int16:
		return v == nil
	case *int32:
		return v == nil
	case *int64:
		return v == nil
	case *uint:
		return v == nil
	case *uint8:
		return v == nil
	case *uint16:
		return v == nil
	case *uint32:
		return v == nil
	case *uint64:
		return v == nil
	case *float32:
		return v == nil
	case *float64:
		return v == nil
	case *string:
		return v == nil
	case *time.Time:
		return v == nil
	case []byte:
		return v == nil
	}
	panic("unsupported type: " + reflect.TypeOf(x).String())
}
