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

func isNumberType(val interface{}) bool {
	if driver.IsValue(val) {
		return false
	}
	valVal := reflect.ValueOf(val)
	switch valVal.Kind() {
	case reflect.Int, reflect.Int64, reflect.Float64:
		return true
	}
	return false
}

func isStringType(val interface{}) (ok bool) {
	if val == nil {
		return false
	}
	_, ok = val.(string)
	return ok
}

type validFn func() (interface{}, bool)

func isValidValue(val interface{}) (interface{}, bool) {
	if fn, ok := val.(validFn); ok {
		return fn()
	}
	return val, true
}

// NoEmpty the value will be skipped if it's empty.
//
// e.g. empty string "", zero 0, nil, nil pointer, empty array/slice.
func NoEmpty(val interface{}) validFn {
	return func() (interface{}, bool) {

		r := reflect.ValueOf(val)

		//check for pointer type
		if r.Kind() == reflect.Ptr {
			if r.IsNil() {
				return val, false
			}
			val = r.Elem().Interface()
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

func removeInvalidValue(values map[string]interface{}) {
	for k, v := range values {
		origVal, ok := isValidValue(v)
		if !ok {
			delete(values, k)
		} else {
			// Set the original value
			values[k] = origVal
		}
	}
}
