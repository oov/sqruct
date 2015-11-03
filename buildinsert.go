package sqruct

// BuildInsertQuery builds SQL such as "INSERT INTO table (column1, column2) VALUES (:column1, NULL)".
func BuildInsertQuery(table string, columns []string, useStructValue []bool) string {
	if len(columns) != len(useStructValue) {
		panic("whoooo")
	}

	// calculate max length
	l := len("INSERT INTO  () VALUES ()") + len(table)
	for i, c := range columns {
		// for columns
		if i != 0 {
			l += 2 // len(", ")
		}
		l += len(c)

		// for values
		if i != 0 {
			l += 2 // len(", ")
		}
		if useStructValue[i] {
			l += len(c) + 1 // len(":"+c)
		} else {
			l += 4 // len("NULL")
		}
	}

	q := make([]byte, 0, l)
	q = append(q, "INSERT INTO "...)
	q = append(q, table...)
	q = append(q, " ("...)
	for i, c := range columns {
		if i != 0 {
			q = append(q, ", "...)
		}
		q = append(q, c...)
	}
	q = append(q, ") VALUES ("...)
	for i, c := range columns {
		if i != 0 {
			q = append(q, ", "...)
		}
		if useStructValue[i] {
			q = append(q, ':')
			q = append(q, c...)
		} else {
			q = append(q, "NULL"...)
		}
	}
	q = append(q, ')')

	return string(q)
}
