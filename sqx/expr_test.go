package sqx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatExpr(t *testing.T) {
	b := ConcatExpr("COALESCE(name,", Expr("CONCAT(?,' ',?)", "f", "l"), ")")
	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "COALESCE(name,CONCAT(?,' ',?))"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"f", "l"}
	assert.Equal(t, expectedArgs, args)
}

func TestConcatExprBadType(t *testing.T) {
	b := ConcatExpr("prefix", 123, "suffix")
	_, _, err := b.ToSql()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "123 is not")
}

func TestExprEscaped(t *testing.T) {
	b := Expr("count(??)", Expr("x"))
	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "count(??)"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{Expr("x")}
	assert.Equal(t, expectedArgs, args)
}

func TestExprRecursion(t *testing.T) {
	{
		b := Expr("count(?)", Expr("nullif(a,?)", "b"))
		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "count(nullif(a,?))"
		assert.Equal(t, expectedSql, sql)

		expectedArgs := []interface{}{"b"}
		assert.Equal(t, expectedArgs, args)
	}
	{
		b := Expr("extract(? from ?)", Expr("epoch"), "2001-02-03")
		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "extract(epoch from ?)"
		assert.Equal(t, expectedSql, sql)

		expectedArgs := []interface{}{"2001-02-03"}
		assert.Equal(t, expectedArgs, args)
	}
	{
		b := Expr("JOIN t1 ON ?", AndP{Eq{"id": 1}, Expr("NOT c1"), Expr("? @@ ?", "x", "y")})
		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "JOIN t1 ON (id = ? AND NOT c1 AND ? @@ ?)"
		assert.Equal(t, expectedSql, sql)

		expectedArgs := []interface{}{1, "x", "y"}
		assert.Equal(t, expectedArgs, args)
	}
}
