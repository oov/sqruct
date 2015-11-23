package sqruct

// PlaceholderGenerator generates placeholders for SQL statements.
type PlaceholderGenerator interface {
	// Placeholder generates next placeholder.
	Placeholder() string
	// Len returns length of placeholder at given index.
	Len(index int) int
}

type genericPlaceholderGenerator struct{}

func (g genericPlaceholderGenerator) Placeholder() string { return "?" }
func (g genericPlaceholderGenerator) Len(int) int         { return 1 }

type postgresPlaceholderGenerator struct {
	c int
}

func (g *postgresPlaceholderGenerator) Placeholder() string {
	g.c++

	var buf [8]byte
	x := g.c
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

func (g postgresPlaceholderGenerator) Len(index int) int {
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
