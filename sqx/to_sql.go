package sqx

import (
	"bytes"
	"io"
)

func ToSql(parts ...interface{}) (sqlStr string, args []interface{}, err error) {

	return toSql("", sqlizerd(parts)...)
}

func sqlizerd(parts []interface{}) (s []Sqlizer) {
	for _, v := range parts {
		switch v := v.(type) {
		case string:
			s = append(s, Expr(v))
		case Sqlizer:
			s = append(s, v)
		}
	}
	return
}

// ToSql wrapper to build sql query, if "sep" or separator is given then each parts going to be separated by the "sep" string
func toSql(sep string, parts ...Sqlizer) (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}

	args, err = appendToSql(parts, sql, sep, args)
	if err != nil {
		return
	}

	sqlStr = sql.String()

	defer logDebug(sqlStr, args, err)
	return
}

func appendToSql(parts []Sqlizer, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		var partSql string
		var partArgs []interface{}
		var err error

		partSql, partArgs, err = p.ToSql()
		if err != nil {
			return nil, err
		} else if len(partSql) == 0 {
			continue
		}

		if i > 0 && len(sep) > 0 {
			_, err := io.WriteString(w, sep)
			if err != nil {
				return nil, err
			}
		} else {
			io.WriteString(w, " ")
		}

		_, err = io.WriteString(w, partSql)
		if err != nil {
			return nil, err
		}
		if len(partArgs) > 0 {
			args = append(args, partArgs...)
		}

	}
	return args, nil
}
