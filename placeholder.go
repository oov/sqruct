package sqruct

import "strings"

// Placeholder generates placeholders for SQL statements.
type Placeholder interface {
	// Next generates next placeholder.
	Next() string
	// Len returns length of placeholder at given index.
	Len(index int) int
	// Rebind replaces from '?' to default placeholder.
	// Every kind of question mark is a replace target because this doesn't interpret a SQL statements.
	Rebind(string) string
}

type genericPlaceholder struct{}

func (*genericPlaceholder) Next() string           { return "?" }
func (*genericPlaceholder) Len(int) int            { return 1 }
func (*genericPlaceholder) Rebind(s string) string { return s }

type postgresPlaceholder struct {
	c int
}

func (ph *postgresPlaceholder) Next() string {
	ph.c++

	var buf [8]byte
	x := ph.c
	i := len(buf) - 1
	for x > 9 {
		buf[i] = byte(x%10 + '0')
		x /= 10
		i--
	}
	buf[i] = byte(x + '0')
	i--
	buf[i] = '$'
	return string(buf[i:])
}

func (ph *postgresPlaceholder) Len(index int) int {
	if index < 0 {
		panic("logic error")
	}
	if index < 9 {
		return 2
	}
	if index < 99 {
		return 3
	}
	if index < 999 {
		return 4
	}
	if index < 9999 {
		return 5
	}
	if index < 99999 {
		return 6
	}
	if index < 999999 {
		return 7
	}
	return 8
}

func (ph *postgresPlaceholder) Rebind(s string) string {
	r := make([]byte, 0, len(s)+8)
	var p int
	for {
		p = strings.IndexByte(s, '?')
		if p == -1 {
			r = append(r, s...)
			break
		}
		r = append(r, s[:p]...)
		r = append(r, ph.Next()...)
		s = s[p+1:]
	}
	return string(r)
}
