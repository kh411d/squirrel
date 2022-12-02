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

type skipperFn func(interface{}) bool

// allow any number including zero 0
var SKIP_NUMBER skipperFn = func(val interface{}) bool { return isNumberType(val) }

// allow nil value
var SKIP_NIL skipperFn = func(val interface{}) bool { return val == nil }

// allow any string including ""
var SKIP_STRING skipperFn = func(val interface{}) bool { return isStringType(val) }

// NoEmptyValue will exclude any empty value ("", 0, nil, ptr type, etc) from the list.
//
// Use skipper arguments SKIP_NUMBER, SKIP_NIL, SKIP_STRING to allow 0, nil , or empty string.
//
// e.g. NoEmptyValue(s, SKIP_NUMBER, SKIP_NIL)
func NoEmptyValue(s map[string]interface{}, skipper ...skipperFn) {
	skipFn := func(val interface{}) (ok bool) {
		for _, fn := range skipper {
			if fn(val) {
				return true
			}
		}
		return
	}

	for k, v := range s {
		r := reflect.ValueOf(v)
		if len(skipper) != 0 && skipFn(v) {
			continue
		}

		if r.Kind() == reflect.Ptr {
			if r.IsNil() {
				v = nil
			} else {
				v = r.Elem().Interface()
				r = reflect.ValueOf(v)
			}
		}

		if !r.IsValid() || v == nil || r.IsZero() {
			delete(s, k)
		}
	}
}
