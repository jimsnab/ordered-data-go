package ordereddata

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	OrderedMap[K comparable, V any] interface {
		Set(k K, v V) (isNew bool)
		Replace(k K, v V)
		Get(k K) (v V, exists bool)
		GetAt(index int) (v V, exists bool)
		MustGet(k K) (v V)
		MustGetString(k K) string
		MustGetInt(k K) int
		MustGetIntOrDefault(k K, defaultInt int) int
		MustGetInt64(k K) int64
		MustGetFloat64(k K) float64
		MustGetFloat64Ptr(k K) (f *float64)
		MustGetBool(k K) bool
		MustGetAt(index int) (v V)
		Map() map[K]V
		Keys() []K
		Len() int
		Values() []V
		Clone() OrderedMap[K, V]
		Copy(other OrderedMap[K, V], keys ...K)
		Remove(k ...K)
		String() string

		encoding.TextUnmarshaler
		json.Marshaler
		json.Unmarshaler
	}

	orderedMap[K comparable, V any] struct {
		m     map[K]V
		order []K
	}

	StringMap interface {
		Set(k string, v any) (isNew bool)
		Replace(k string, v any)
		Get(k string) (v any, exists bool)
		GetAt(index int) (v any, exists bool)
		MustGet(k string) (v any)
		MustGetString(k string) string
		MustGetInt(k string) int
		MustGetIntOrDefault(k string, defaultInt int) int
		MustGetInt64(k string) int64
		MustGetFloat64(k string) float64
		MustGetFloat64Ptr(k string) (f *float64)
		MustGetBool(k string) bool
		MustGetAt(index int) (v any)
		Map() map[string]any
		Keys() []string
		Len() int
		Values() []any
		Clone() OrderedMap[string, any]
		Copy(other OrderedMap[string, any], keys ...string)
		Remove(k ...string)
		String() string

		encoding.TextUnmarshaler
		json.Marshaler
		json.Unmarshaler
	}
)

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return &orderedMap[K, V]{
		m:     map[K]V{},
		order: []K{},
	}
}

func NewOrderedMapN[K comparable, V any](length int) OrderedMap[K, V] {
	return &orderedMap[K, V]{
		m:     make(map[K]V, length),
		order: make([]K, 0, length),
	}
}

func NewStringMap() StringMap {
	return NewOrderedMap[string, any]()
}

func NewStringMapN(length int) StringMap {
	return NewOrderedMapN[string, any](length)
}

func (om *orderedMap[K, V]) Clone() OrderedMap[K, V] {
	clone := NewOrderedMapN[K, V](om.Len())
	for _, k := range om.Keys() {
		v := om.m[k]
		clone.Set(k, v)
	}
	return clone
}

func (om *orderedMap[K, V]) Set(k K, v V) (isNew bool) {
	_, exists := om.m[k]
	if !exists {
		om.order = append(om.order, k)
	} else {
		isNew = true
	}
	om.m[k] = v
	return
}

func (om *orderedMap[K, V]) Replace(k K, v V) {
	_, exists := om.m[k]
	if !exists {
		panic("key doesn't exist")
	}
	om.m[k] = v
}

func (om *orderedMap[K, V]) Get(k K) (v V, exists bool) {
	v, exists = om.m[k]
	return
}

func (om *orderedMap[K, V]) GetAt(index int) (v V, exists bool) {
	if index < 0 || index >= len(om.order) {
		return
	}
	v, exists = om.m[om.order[index]]
	return
}

func (om *orderedMap[K, V]) Copy(other OrderedMap[K, V], keys ...K) {
	for _, key := range keys {
		v, exists := other.Get(key)
		if exists {
			om.Set(key, v)
		} else {
			om.Remove(key)
		}
	}
}

func (om *orderedMap[K, V]) MustGet(k K) V {
	return om.m[k]
}

func (om *orderedMap[K, V]) MustGetString(k K) string {
	vany := om.m[k]
	rv := reflect.ValueOf(vany)
	if !rv.IsValid() {
		return ""
	}

	kind := rv.Kind()
	if (kind == reflect.Interface || kind == reflect.Pointer) && rv.IsNil() {
		return ""
	}

	rv = reflect.Indirect(rv)
	kind = rv.Kind()
	if (kind == reflect.Interface || kind == reflect.Pointer) && rv.IsNil() {
		return ""
	}

	return fmt.Sprintf("%v", rv)
}

func (om *orderedMap[K, V]) MustGetFloat64Ptr(k K) *float64 {
	str := om.MustGetString(k)
	if str == "" {
		return nil
	}
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}
	return &v
}

func (om *orderedMap[K, V]) MustGetInt(k K) int {
	str := om.MustGetString(k)
	v, _ := strconv.Atoi(str)
	return v
}

func (om *orderedMap[K, V]) MustGetIntOrDefault(k K, defaultInt int) int {
	str := om.MustGetString(k)
	v, err := strconv.Atoi(str)
	if err != nil {
		return defaultInt
	}
	return v
}

func (om *orderedMap[K, V]) MustGetInt64(k K) int64 {
	str := om.MustGetString(k)
	v, _ := strconv.ParseInt(str, 10, 64)
	return v
}

func (om *orderedMap[K, V]) MustGetFloat64(k K) float64 {
	str := om.MustGetString(k)
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func (om *orderedMap[K, V]) MustGetBool(k K) bool {
	str := om.MustGetString(k)
	if str == "Y" {
		return true
	}
	if str == "N" {
		return false
	}
	v, _ := strconv.ParseBool(str)
	return v
}

func (om *orderedMap[K, V]) MustGetAt(index int) (v V) {
	if index < 0 || index >= len(om.order) {
		return
	}
	return om.m[om.order[index]]
}

func (om *orderedMap[K, V]) Map() map[K]V {
	return om.m
}

func (om *orderedMap[K, V]) Keys() []K {
	return om.order
}

func (om *orderedMap[K, V]) Len() int {
	return len(om.order)
}

func (om *orderedMap[K, V]) Values() []V {
	vals := make([]V, 0, len(om.order))
	for _, k := range om.order {
		vals = append(vals, om.m[k])
	}
	return vals
}

func (om *orderedMap[K, V]) MarshalJSON() (text []byte, err error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	w.WriteRune('{')

	for n, k := range om.order {
		if n != 0 {
			w.WriteRune(',')
		}

		w.WriteRune('"')
		w.WriteString(fmt.Sprintf("%v", k))
		w.WriteRune('"')
		w.WriteRune(':')

		var value []byte
		if value, err = json.Marshal(om.m[k]); err != nil {
			return
		}

		w.Write(value)
	}

	w.WriteRune('}')
	w.Flush()

	text = buf.Bytes()
	return
}

func (om *orderedMap[K, V]) UnmarshalJSON(text []byte) (err error) {
	oj := newOrderedJson(text)

	m, _, err := oj.peekJsonMap(parserPos{0, 1})
	if err != nil {
		return
	}

	om.m = make(map[K]V, m.Len())
	om.order = make([]K, 0, m.Len())

	for _, k := range m.Keys() {
		var key K

		switch any(key).(type) {
		case K:
			key = reflect.ValueOf(k).Interface().(K)
		default:
			a := any(key)
			jm, is := a.(json.Unmarshaler)
			if is {
				if err = jm.UnmarshalJSON([]byte(k)); err != nil {
					return
				}
			} else {
				tm, is := a.(encoding.TextUnmarshaler)
				if is {
					if err = tm.UnmarshalText([]byte(k)); err != nil {
						return
					}
				} else {
					err = errors.New("ordered map key must implement an unmarshaler")
					return
				}
			}
		}

		om.order = append(om.order, key)
		v, _ := m.MustGet(k).(V)
		om.m[key] = v
	}

	return
}

func (om *orderedMap[K, V]) UnmarshalText(text []byte) (err error) {
	oj := newOrderedJson(text)

	m, _, err := oj.peekJsonMap(parserPos{0, 1})
	if err != nil {
		return
	}

	for _, k := range m.Keys() {
		var key K
		tm := any(key).(encoding.TextUnmarshaler)
		if err = tm.UnmarshalText([]byte(k)); err != nil {
			return
		}

		v := m.MustGet(k)
		om.order = append(om.order, key)
		om.m[key] = v.(V)
	}

	return
}

func (om *orderedMap[K, V]) Remove(keys ...K) {
	toRemove := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		toRemove[k] = struct{}{}
	}

	newOrder := make([]K, 0, len(om.order))
	for _, key := range om.order {
		if _, exists := toRemove[key]; exists {
			delete(om.m, key)
		} else {
			newOrder = append(newOrder, key)
		}
	}
	om.order = newOrder
}

func (om orderedMap[K, V]) String() string {
	var sb strings.Builder

	sb.WriteRune('{')

	for _, k := range om.order {
		if sb.Len() > 1 {
			sb.WriteString(", ")
		}
		v := om.m[k]

		rv := reflect.ValueOf(v)
		if !rv.IsValid() {
			sb.WriteString(fmt.Sprintf(`%s: nil`, jsonEscape(k)))
		} else if rv.CanInterface() {
			anyVal := reflect.Indirect(rv).Interface()
			if anyVal == nil {
				sb.WriteString(fmt.Sprintf(`%s: nil`, jsonEscape(k)))
			} else {
				str, is := anyVal.(string)
				if is {
					str = fmt.Sprintf(`"%s"`, str)
				} else {
					str, is = convertToString(anyVal)
				}
				if is {
					sb.WriteString(fmt.Sprintf(`%s: %s`, jsonEscape(k), str))
				} else {
					sb.WriteString(fmt.Sprintf(`%s: %#v`, jsonEscape(k), anyVal))
				}
			}
		} else {
			sb.WriteString(fmt.Sprintf(`%s: %s`, jsonEscape(k), jsonEscape(v)))
		}
	}

	sb.WriteRune('}')

	return sb.String()
}

func jsonEscape[V any](v V) string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // To prevent escaping HTML characters
	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	escapedString := buf.String()

	// Remove the trailing newline added by Encode
	escapedString = escapedString[:len(escapedString)-1]

	return escapedString
}

func convertToString(v any) (string, bool) {
	if str, ok := v.(fmt.Stringer); ok {
		return str.String(), true
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		if str, ok := rv.Interface().(fmt.Stringer); ok {
			return str.String(), true
		}
	}

	return "", false
}
