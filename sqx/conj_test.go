package sqx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyAndToSql(t *testing.T) {
	sql, args, err := And{}.ToSql()
	assert.NoError(t, err)

	expectedSql := "(1=1)"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestEmptyOrToSql(t *testing.T) {
	sql, args, err := Or{}.ToSql()
	assert.NoError(t, err)

	expectedSql := "(1=0)"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestConjWhere(t *testing.T) {
	x := Where{
		Eq{"a": 1, "b": 2, "c": "three"},
		Like{"d": "%kambing"},
	}

	s, args, err := x.ToSql()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s, " WHERE a = ? AND b = ? AND c = ? AND d LIKE ? ")
	assert.Equal(t, args, []interface{}{1, 2, "three", "%kambing"})
}

func TestConjWhereOR(t *testing.T) {
	x := Where{
		Or{
			Eq{"a": 1, "b": 2, "c": "three"},
			Like{"d": "%kambing"},
		},
	}

	s, args, err := x.ToSql()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s, " WHERE a = ? AND b = ? AND c = ? OR d LIKE ? ")
	assert.Equal(t, args, []interface{}{1, 2, "three", "%kambing"})
}

func TestConjWhereEmptyValue(t *testing.T) {
	x := Where{
		Eq{},
		Like{},
	}

	s, args, err := x.ToSql()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s, " WHERE (1=1) ")
	assert.Equal(t, args, []interface{}(nil))
}

func TestConjAND(t *testing.T) {
	x := And{
		Eq{"a": 1, "b": 2, "c": "three"},
		Like{"d": "%kambing"},
	}

	s, args, err := x.ToSql()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s, " a = ? AND b = ? AND c = ? AND d LIKE ? ")
	assert.Equal(t, args, []interface{}{1, 2, "three", "%kambing"})
}

func TestConjAndP(t *testing.T) {
	x := AndP{
		Eq{"a": 1, "b": 2, "c": "three"},
		Like{"d": "%kambing"},
	}

	s, args, err := x.ToSql()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, s, "(a = ? AND b = ? AND c = ? AND d LIKE ?)")
	assert.Equal(t, args, []interface{}{1, 2, "three", "%kambing"})
}
