package sqruct

// Columns returns string such as "prefix.column0, prefix.column1, prefix.column2"
func Columns(prefix string, columns []string) string {
	if len(columns) == 0 {
		return ""
	}

	if prefix != "" {
		prefix += "."
	}
	// +2 = len(", ")
	// -2 = strip last ", "
	l := (len(prefix)+2)*len(columns) - 2
	for _, v := range columns {
		l += len(v)
	}
	buf := make([]byte, 0, l)
	for i, v := range columns {
		if i != 0 {
			buf = append(buf, ", "...)
		}
		buf = append(buf, prefix...)
		buf = append(buf, v...)
	}
	return string(buf)
}
