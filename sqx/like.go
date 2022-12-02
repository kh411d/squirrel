package sqx

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Like is syntactic sugar for use with LIKE conditions.
// Ex:
//     .Where(Like{"name": "%irrel"})
type Like map[string]interface{}

func (lk Like) toSql(opr string) (sql string, args []interface{}, err error) {
	var exprs []string
	for key, val := range lk {
		expr := ""

		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}

		if val == nil {
			err = fmt.Errorf("cannot use null with like operators")
			return
		} else {
			if isListType(val) {
				err = fmt.Errorf("cannot use array or slice with like operators")
				return
			} else {
				expr = fmt.Sprintf("%s %s ?", key, opr)
				args = append(args, val)
			}
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

func (lk Like) ToSql() (sql string, args []interface{}, err error) {
	return lk.toSql("LIKE")
}

// NotLike is syntactic sugar for use with LIKE conditions.
// Ex:
//     .Where(NotLike{"name": "%irrel"})
type NotLike Like

func (nlk NotLike) ToSql() (sql string, args []interface{}, err error) {
	return Like(nlk).toSql("NOT LIKE")
}

// ILike is syntactic sugar for use with ILIKE conditions.
// Ex:
//    .Where(ILike{"name": "sq%"})
type ILike Like

func (ilk ILike) ToSql() (sql string, args []interface{}, err error) {
	return Like(ilk).toSql("ILIKE")
}

// NotILike is syntactic sugar for use with ILIKE conditions.
// Ex:
//    .Where(NotILike{"name": "sq%"})
type NotILike Like

func (nilk NotILike) ToSql() (sql string, args []interface{}, err error) {
	return Like(nilk).toSql("NOT ILIKE")
}
