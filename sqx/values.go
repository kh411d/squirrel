package sqx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Values generate VALUES for insert statement
type Values [][]interface{}

func (insv Values) ToSql() (sql string, args []interface{}, err error) {
	b := &bytes.Buffer{}

	args, err = appendValuesToSQL(b, insv, args)
	if err != nil {
		return
	}

	sql = b.String()
	return
}

func appendValuesToSQL(w io.Writer, insValues Values, args []interface{}) ([]interface{}, error) {
	if len(insValues) == 0 {
		return args, errors.New("values for insert statements are not set")
	}

	io.WriteString(w, "VALUES ")

	valuesStrings := make([]string, len(insValues))
	for r, row := range insValues {
		valueStrings := make([]string, len(row))
		for v, val := range row {
			if vs, ok := val.(Sqlizer); ok {
				vsql, vargs, err := vs.ToSql()
				if err != nil {
					return nil, err
				}
				valueStrings[v] = vsql
				args = append(args, vargs...)
			} else {
				valueStrings[v] = "?"
				args = append(args, val)
			}
		}
		valuesStrings[r] = fmt.Sprintf("(%s)", strings.Join(valueStrings, ","))
	}

	io.WriteString(w, strings.Join(valuesStrings, ","))

	return args, nil
}
