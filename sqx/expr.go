package sqx

import (
	"bytes"
	"fmt"
	"strings"
)

type expr struct {
	sql  string
	args []interface{}
}

// Expr builds an expression from a SQL fragment and arguments.
//
// Ex:
//     Expr("FROM_UNIXTIME(?)", t)
func Expr(sql string, args ...interface{}) Sqlizer {
	return expr{sql: sql, args: args}
}

func (e expr) ToSql() (sql string, args []interface{}, err error) {
	simple := true
	for _, arg := range e.args {
		if _, ok := arg.(Sqlizer); ok {
			simple = false
		}
	}
	if simple {
		return e.sql, e.args, nil
	}

	buf := &bytes.Buffer{}
	ap := e.args
	sp := e.sql

	var isql string
	var iargs []interface{}

	for err == nil && len(ap) > 0 && len(sp) > 0 {
		i := strings.Index(sp, "?")
		if i < 0 {
			// no more placeholders
			break
		}
		if len(sp) > i+1 && sp[i+1:i+2] == "?" {
			// escaped "??"; append it and step past
			buf.WriteString(sp[:i+2])
			sp = sp[i+2:]
			continue
		}

		if as, ok := ap[0].(Sqlizer); ok {
			// sqlizer argument; expand it and append the result
			isql, iargs, err = as.ToSql()
			buf.WriteString(sp[:i])
			buf.WriteString(isql)
			args = append(args, iargs...)
		} else {
			// normal argument; append it and the placeholder
			buf.WriteString(sp[:i+1])
			args = append(args, ap[0])
		}

		// step past the argument and placeholder
		ap = ap[1:]
		sp = sp[i+1:]
	}

	// append the remaining sql and arguments
	buf.WriteString(sp)
	return buf.String(), append(args, ap...), err
}

type concatExpr []interface{}

func (ce concatExpr) ToSql() (sql string, args []interface{}, err error) {
	for _, part := range ce {
		switch p := part.(type) {
		case string:
			sql += p
		case Sqlizer:
			pSql, pArgs, err := p.ToSql()
			if err != nil {
				return "", nil, err
			}
			sql += pSql
			args = append(args, pArgs...)
		default:
			return "", nil, fmt.Errorf("%#v is not a string or Sqlizer", part)
		}
	}
	return
}

// ConcatExpr builds an expression by concatenating strings and other expressions.
//
// Ex:
//     name_expr := Expr("CONCAT(?, ' ', ?)", firstName, lastName)
//     ConcatExpr("COALESCE(full_name,", name_expr, ")")
func ConcatExpr(parts ...interface{}) concatExpr {
	return concatExpr(parts)
}
