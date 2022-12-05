package sqx

import (
	"database/sql/driver"
	"reflect"
	"sort"
)

const (
	// Portable true/false literals.
	sqlTrue  = "(1=1)"
	sqlFalse = "(1=0)"
	sqlEmpty = " "
)

// Sqlizer is the interface that wraps the ToSql method.
//
// ToSql returns a SQL representation of the Sqlizer, along with a slice of args
// as passed to e.g. database/sql.Exec. It can also return an error.
type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

func getSortedKeys(exp map[string]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func isListType(val interface{}) bool {
	if driver.IsValue(val) {
		return false
	}
	valVal := reflect.ValueOf(val)
	return valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice
}

type validFn func(skip bool) (interface{}, bool)

func isValidValue(val interface{}, skip bool) (interface{}, bool) {

	if fn, ok := val.(validFn); ok {
		return fn(skip)
	}
	return val, true
}

/*
Statement clause should void if value is empty.

e.g. empty string "", zero 0, nil, nil pointer, empty array/slice.

	val := ""
	Eq{"table_a.name": val, "table_b.address":"jakarta"}
	=> `table_a.name = "" AND table_b.address = "jakarta"`

	Eq{"table_a.name": sqx.NoEmpty(val), "table_b.address":"jakarta"}
	=> `table_b.address = "jakarta"` //table_a.name clause is skipped
*/
func NoEmpty(val interface{}) validFn {
	return func(skip bool) (interface{}, bool) {

		if skip {
			// return driver value type only
			return val, true
		}

		r := reflect.ValueOf(val)

		//check for pointer type
		if r.Kind() == reflect.Ptr {
			if r.IsNil() {
				return nil, false
			}
			val = r.Elem().Interface()
			//get the underlying value
			r = reflect.ValueOf(val)
		}
		//check for empty array and slice
		if (r.Kind() == reflect.Array || r.Kind() == reflect.Slice) && r.Len() == 0 {
			return val, false
		}

		if val == nil {
			return val, false
		}

		return val, !r.IsZero()
	}
}
