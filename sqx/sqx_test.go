package sqx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEmptyValue(t *testing.T) {
	var i interface{}
	var intVal int
	x := Eq{
		"a": NoEmpty(""),
		"b": NoEmpty(i),
		"c": NoEmpty(nil),
		"d": NoEmpty(&intVal),
		"e": NoEmpty(0),
		"f": NoEmpty(12345),
	}

	sql, args, _ := x.ToSql()

	assert.Equal(t, sql, "f = ?")
	assert.Equal(t, args, []interface{}{12345})
}

func TestNoEmptyValueWitPtr(t *testing.T) {

	var nilIntVal int
	var nilInterface interface{}
	var intVal int = 0
	s := []string{"abc"}
	n := struct {
		Name string
	}{
		Name: "name",
	}
	g := []string{}
	x := Eq{
		"a": &intVal,
		"b": &s,
		"c": &n,
		"d": NoEmpty(&nilIntVal),    // should be excluded
		"e": NoEmpty(nilInterface),  // should be excluded
		"f": NoEmpty(&nilInterface), // should be excluded
		"g": NoEmpty(g),
	}

	sql, args, _ := x.ToSql()

	assert.Equal(t, sql, "a = ? AND b IN (?) AND c = ?")
	assert.Equal(t, args, []interface{}{0, "abc", struct{ Name string }{Name: "name"}})
}
