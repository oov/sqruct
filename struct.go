package sqruct

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func extractStructFieldType(structField string) (string, error) {
	const pre = "package p\ntype T struct{\n"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", pre+structField+"\n}", 0)
	if err != nil {
		return "", err
	}
	var pos, end int
	ast.Inspect(f.Decls[0], func(n ast.Node) bool {
		st, ok := n.(*ast.StructType)
		if !ok {
			return true
		}
		f := st.Fields.List[0]
		pos, end = int(f.Type.Pos()), int(f.Type.End())
		return false
	})
	if pos == 0 {
		return "", fmt.Errorf(`could not parse as struct field %q`, structField)
	}
	return structField[pos-len(pre)-1 : end-len(pre)-1], nil
}

func parseStructTag(structField string) (pos int, end int, err error) {
	const pre = "package p\ntype T struct{\n"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", pre+structField+"\n}", 0)
	if err != nil {
		return 0, 0, err
	}
	ast.Inspect(f.Decls[0], func(n ast.Node) bool {
		st, ok := n.(*ast.StructType)
		if !ok {
			return true
		}
		f := st.Fields.List[0]
		if f.Tag != nil {
			pos, end = int(f.Tag.Pos()), int(f.Tag.End())
		} else {
			pos, end = int(f.End()), int(f.End())
		}
		return false
	})
	if pos == 0 {
		return 0, 0, fmt.Errorf(`could not parse as struct field %q`, structField)
	}
	return pos - len(pre) - 1, end - len(pre) - 1, nil
}

func replaceStructTag(structField string, replacer func(s string) (string, error)) (string, error) {
	pos, end, err := parseStructTag(structField)
	if err != nil {
		return "", err
	}
	tag := structField[pos:end]
	if tag != "" {
		tag, err = strconv.Unquote(tag)
		if err != nil {
			return "", err
		}
	}

	tag, err = replacer(tag)
	if err != nil {
		return "", err
	}

	if strconv.CanBackquote(tag) {
		tag = "`" + tag + "`"
	} else {
		tag = strconv.Quote(tag)
	}

	if pos == end {
		tag = " " + tag
	}
	return structField[:pos] + tag + structField[end:], nil
}
