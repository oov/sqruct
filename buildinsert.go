package sqruct

// buildInsert builds SQL such as "INSERT INTO table (column1, column2) VALUES (:column1, :column2)".
func buildInsert(table string, columns []string, autoIncrCol int, defValue string, useReturning bool) string {
	// calculate max length
	l := len("INSERT INTO  () VALUES ()") + len(table)
	for i, c := range columns {
		if i != 0 {
			l += 4 // len(", ") for columns + values
		}
		if i == autoIncrCol {
			l += len(c) + len(defValue) // len(c) for columns + len(defValue) for values
			if useReturning {
				l += 11 + len(c) // len(" RETURNING "+c)
			}
		} else {
			l += len(c)*2 + 1 // len(c) for columns + len(":"+c) for values
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
		if i == autoIncrCol {
			q = append(q, defValue...)
		} else {
			q = append(q, ':')
			q = append(q, c...)
		}
	}
	q = append(q, ')')
	if useReturning && autoIncrCol >= 0 && autoIncrCol < len(columns) {
		q = append(q, " RETURNING "...)
		q = append(q, columns[autoIncrCol]...)
	}

	return string(q)
}
