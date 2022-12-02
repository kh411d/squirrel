package sqx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLikeToSql(t *testing.T) {
	b := Like{"name": "%irrel"}
	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "name LIKE ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"%irrel"}
	assert.Equal(t, expectedArgs, args)
}

func TestNotLikeToSql(t *testing.T) {
	b := NotLike{"name": "%irrel"}
	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "name NOT LIKE ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"%irrel"}
	assert.Equal(t, expectedArgs, args)
}

func TestILikeToSql(t *testing.T) {
	b := ILike{"name": "sq%"}
	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "name ILIKE ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"sq%"}
	assert.Equal(t, expectedArgs, args)
}

func TestNotILikeToSql(t *testing.T) {
	b := NotILike{"name": "sq%"}
	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "name NOT ILIKE ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"sq%"}
	assert.Equal(t, expectedArgs, args)
}
