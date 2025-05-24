// Package ptrmap provides a function to collect pointers to struct fields from a given struct.
package ptrmap

import (
	"maps"
	"reflect"
	"strings"
)

// ExcludePredicate is a function type that defines a predicate for excluding fields from collection.
type ExcludePredicate func(field reflect.StructField, value reflect.Value) bool

// Collector is a struct that collects pointers to fields of a struct, excluding fields based on the provided
// predicates.
type Collector struct {
	excludePredicates []ExcludePredicate
}

// New returns a new Collector with the given exclude predicates.
func New(predicates ...ExcludePredicate) *Collector {
	return &Collector{
		excludePredicates: predicates,
	}
}

// DefaultCollector returns a Collector with default exclusion rules.
func DefaultCollector() *Collector {
	return New(
		func(ft reflect.StructField, _ reflect.Value) bool {
			return ft.PkgPath != ""
		},
		func(_ reflect.StructField, fv reflect.Value) bool {
			return fv.Kind() != reflect.Struct
		},
		func(ft reflect.StructField, _ reflect.Value) bool {
			tag := ft.Tag.Get("mapstructure")

			return strings.HasPrefix(tag, ",") || tag == "-"
		},
	)
}

// AddExcludePredicate adds an exclusion predicate to the Collector.
func (c *Collector) AddExcludePredicate(pred ExcludePredicate) {
	c.excludePredicates = append(c.excludePredicates, pred)
}

// Collect recursively collects pointers to fields of a struct.
func (c *Collector) Collect(v any) map[string]any {
	return c.collect(v, "")
}

// collect is a helper function that performs the actual collection of pointers to struct fields.
func (c *Collector) collect(v any, prefix string) map[string]any {
	out := make(map[string]any)
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("collect expects a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := range typ.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)

		excluded := false

		for _, pred := range c.excludePredicates {
			if pred(fieldType, field) {
				excluded = true

				break
			}
		}

		if excluded {
			continue
		}

		name := fieldType.Tag.Get("mapstructure")
		if name == "" {
			name = strings.ToLower(fieldType.Name)
		}

		if field.CanAddr() {
			key := name
			if prefix != "" {
				key = prefix + "." + name
			}

			out[key] = field.Addr().Interface()
			maps.Copy(out, c.collect(field.Addr().Interface(), key))
		}
	}

	return out
}
