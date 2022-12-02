package sqx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEmptyValue(t *testing.T) {
	var i interface{}
	var intVal int
	x := map[string]interface{}{
		"a": "",
		"b": i,
		"c": nil,
		"d": &intVal,
		"e": 0,
		"f": 12345,
	}

	//should not panic
	NoEmptyValue(x)

	if !reflect.DeepEqual(x, map[string]interface{}{"f": 12345}) {
		t.Error("data missmatch")
	}
}

func TestNoEmptyValueWithSkipper(t *testing.T) {
	var i interface{}
	var intVal int

	x := map[string]interface{}{
		"a": "",
		"b": i,
		"c": nil,
		"d": &intVal,
		"e": 0,
		"f": 12345,
	}

	//should not panic
	NoEmptyValue(x, SKIP_NIL)
	//fmt.Println(x)

	if !reflect.DeepEqual(x, map[string]interface{}{"b": nil, "c": nil, "f": 12345}) {
		t.Error("data missmatch")
	}

	y := map[string]interface{}{
		"a": "",
		"b": i,
		"c": nil,
		"d": &intVal,
		"e": 0,
		"f": 12345,
	}

	//should not panic
	NoEmptyValue(y, SKIP_NIL, SKIP_NUMBER)
	//fmt.Println(y)

	if !reflect.DeepEqual(y, map[string]interface{}{"b": nil, "c": nil, "e": 0, "f": 12345}) {
		t.Error("data missmatch")
	}

	z := map[string]interface{}{
		"a": "",
		"b": i,
		"c": nil,
		"d": &intVal,
		"e": 0,
		"f": 12345,
	}

	//should not panic
	NoEmptyValue(z, SKIP_NIL, SKIP_NUMBER, SKIP_STRING)
	//fmt.Println(z)

	if !reflect.DeepEqual(z, map[string]interface{}{"a": "", "b": nil, "c": nil, "e": 0, "f": 12345}) {
		t.Error("data missmatch")
	}
}

func TestNoEmptyValueWitPtr(t *testing.T) {

	var nilIntVal int
	var nilInterface interface{}
	var intVal int = 1234
	s := []string{"abc"}
	n := struct {
		Name string
	}{
		Name: "name",
	}
	x := map[string]interface{}{
		"a": &intVal,
		"b": &s,
		"c": &n,
		"d": &nilIntVal,    // should be excluded
		"e": nilInterface,  // should be excluded
		"f": &nilInterface, // should be excluded
	}

	//should not panic
	NoEmptyValue(x)
	assert.Equal(t, x, map[string]interface{}{"a": &intVal, "b": &s, "c": &n})
}
