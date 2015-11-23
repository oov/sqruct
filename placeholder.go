package sqruct

import "strconv"

type PlaceholderGenerator interface {
	Placeholder() string
	Len(index int) int
}

type genericPlaceholderGenerator struct{}

func (g genericPlaceholderGenerator) Placeholder() string { return "?" }
func (g genericPlaceholderGenerator) Len(int) int         { return 1 }

type postgresPlaceholderGenerator int

func (g *postgresPlaceholderGenerator) Placeholder() string {
	*g++
	return "$" + strconv.Itoa(int(*g))
}

func (g postgresPlaceholderGenerator) Len(index int) int {
	switch {
	case index < 9:
		return 2
	case index < 99:
		return 3
	case index < 999:
		return 4
	}
	return 5
}
