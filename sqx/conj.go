package sqx

import (
	"fmt"
	"strings"
)

type conj []Sqlizer

var defaultConjFormat string = "(%s)"
var noConjFormat string = " %s "

func (c conj) join(sep, defaultExpr string, defaultFormat string) (sql string, args []interface{}, err error) {
	if len(c) == 0 {
		return defaultExpr, []interface{}{}, nil
	}
	var sqlParts []string
	for _, sqlizer := range c {
		partSQL, partArgs, err := sqlizer.ToSql()
		if err != nil {
			return "", nil, err
		}
		if partSQL != "" {
			sqlParts = append(sqlParts, partSQL)
			args = append(args, partArgs...)
		}
	}
	if len(sqlParts) > 0 {
		sql = fmt.Sprintf(defaultFormat, strings.Join(sqlParts, sep))
	}
	return
}

/*
And conjunction Sqlizers
	x := And{
		Eq{"a": 1, "b": 2, "c": "three"},
		Like{"d": "%kambing"},
	}

	s, args, err := x.ToSql()

	assert.Equal(t, s, " a = ? AND b = ? AND c = ? AND d LIKE ? ")
	assert.Equal(t, args, []interface{}{1, 2, "three", "%kambing"})
*/
type And conj

func (a And) ToSql() (string, []interface{}, error) {
	return conj(a).join(" AND ", sqlTrue, noConjFormat)
}

// Or conjunction Sqlizers
type Or conj

func (o Or) ToSql() (string, []interface{}, error) {
	return conj(o).join(" OR ", sqlFalse, noConjFormat)
}

// AndP conjunction Sqlizers with bracket
type AndP conj

func (a AndP) ToSql() (string, []interface{}, error) {
	return conj(a).join(" AND ", sqlTrue, defaultConjFormat)
}

// OrP conjunction Sqlizers with bracket
type OrP conj

func (o OrP) ToSql() (string, []interface{}, error) {
	return conj(o).join(" OR ", sqlFalse, defaultConjFormat)
}

/*
Where conjunction:
	x := Where{
		Eq{"a": 1, "b": 2, "c": "three"},
		Like{"d": "%kambing"},
	}

	s, args, err := x.ToSql()
	s    => " WHERE a = ? AND b = ? AND c = ? AND d LIKE ? "
	args => []interface{}{1, 2, "three", "%kambing"}
*/
type Where conj

func (w Where) ToSql() (s string, args []interface{}, err error) {
	s, args, err = conj(w).join(" AND ", sqlEmpty, noConjFormat)
	if len(s) > 1 {
		s = fmt.Sprintf(" WHERE %s ", strings.TrimSpace(s))
	}
	return
}

/*
Having conjunction:
	x := Having{
		Eq{"a": 1, "b": 2, "c": "three"},
		Like{"d": "%kambing"},
	}

	s, args, err := x.ToSql()
	s    => " Having a = ? AND b = ? AND c = ? AND d LIKE ? "
	args => []interface{}{1, 2, "three", "%kambing"}
*/
type Having conj

func (h Having) ToSql() (s string, args []interface{}, err error) {
	s, args, err = conj(h).join(" AND ", sqlEmpty, noConjFormat)
	if len(s) > 1 {
		s = fmt.Sprintf(" HAVING %s ", strings.TrimSpace(s))
	}
	return
}
