package sqx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// See to_sql_test.go for insert statement implementation
func TestValuesSingle(t *testing.T) {

	v := Values{
		{"a", "b", "c"},
	}

	sql, args, err := v.ToSql()

	assert.NoError(t, err)
	assert.Equal(t, sql, "VALUES (?,?,?)")
	assert.Equal(t, args, []interface{}{"a", "b", "c"})
}

func TestValuesMultiple(t *testing.T) {

	data := Values{
		{"a", "b", "c"},
		{"aa", "bb", "cc"},
		{"aaa", "bbb", "ccc"},
	}

	sql, args, err := data.ToSql()

	assert.NoError(t, err)
	assert.Equal(t, sql, "VALUES (?,?,?),(?,?,?),(?,?,?)")
	assert.Equal(t, args, []interface{}{"a", "b", "c", "aa", "bb", "cc", "aaa", "bbb", "ccc"})
}
