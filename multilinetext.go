package sqruct

import "strings"

// MultiLineText represents multi line text.
type MultiLineText string

// AddPrefix adds prefix every lines.
func (m MultiLineText) AddPrefix(p string) MultiLineText {
	return MultiLineText(p + strings.NewReplacer("\r\n", "\r\n"+p, "\n", "\n"+p).Replace(string(m)))
}

// AddPostfix adds postfix every lines.
func (m MultiLineText) AddPostfix(p string) MultiLineText {
	return MultiLineText(strings.NewReplacer("\r\n", p+"\r\n", "\n", p+"\n").Replace(string(m)) + p)
}

// String implements fmt.Stringer interface.
func (m MultiLineText) String() string {
	return string(m)
}
