package safesimplejson

import (
	"github.com/bitly/go-simplejson"
	"io"
	"sync"
)

func Version() string {
	return "0.5.0"
}

type Json struct {
	l  *sync.Mutex
	sj *simplejson.Json
}

// NewJson returns a pointer to a new `Json` object
// after unmarshaling `body` bytes
func NewJson(body []byte) (*Json, error) {
	j, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	return &Json{new(sync.Mutex), j}, nil
}

// New returns a pointer to a new, empty `Json` object
func New() *Json {
	return &Json{new(sync.Mutex), simplejson.New()}
}

func NewFromReader(r io.Reader) (*Json, error) {
	j, err := simplejson.NewFromReader(r)
	if err != nil {
		return nil, err
	}
	return &Json{new(sync.Mutex), j}, nil
}

// Implements the json.Unmarshaler interface.
func (j *Json) UnmarshalJSON(p []byte) error {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.UnmarshalJSON(p)
}

// Interface returns the underlying data
func (j *Json) Interface() interface{} {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Interface()
}

// Encode returns its marshaled data as `[]byte`
func (j *Json) Encode() ([]byte, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MarshalJSON()
}

// EncodePretty returns its marshaled data as `[]byte` with indentation
func (j *Json) EncodePretty() ([]byte, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.EncodePretty()
}

// Implements the json.Marshaler interface.
func (j *Json) MarshalJSON() ([]byte, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MarshalJSON()
}

// Set modifies `Json` map by `key` and `value`
// Useful for changing single key/value in a `Json` object easily.
func (j *Json) Set(key string, val interface{}) {
	j.l.Lock()
	defer j.l.Unlock()
	j.sj.Set(key, val)
}

// SetPath modifies `Json`, recursively checking/creating map keys for the supplied path,
// and then finally writing in the value
func (j *Json) SetPath(branch []string, val interface{}) {
	j.l.Lock()
	defer j.l.Unlock()
	j.sj.SetPath(branch, val)
}

// Del modifies `Json` map by deleting `key` if it is present.
func (j *Json) Del(key string) {
	j.l.Lock()
	defer j.l.Unlock()
	j.sj.Del(key)
}

// Get returns a pointer to a new `Json` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *Json) Get(key string) *Json {
	j.l.Lock()
	defer j.l.Unlock()
	return &Json{new(sync.Mutex), j.sj.Get(key)}
}

// GetPath searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetPath("top_level", "dict")
func (j *Json) GetPath(branch ...string) *Json {
	j.l.Lock()
	defer j.l.Unlock()
	return &Json{new(sync.Mutex), j.sj.GetPath(branch...)}
}

// GetIndex returns a pointer to a new `Json` object
// for `index` in its `array` representation
//
// this is the analog to Get when accessing elements of
// a json array instead of a json object:
//    js.Get("top_level").Get("array").GetIndex(1).Get("key").Int()
func (j *Json) GetIndex(index int) *Json {
	j.l.Lock()
	defer j.l.Unlock()
	return &Json{new(sync.Mutex), j.sj.GetIndex(index)}
}

// CheckGet returns a pointer to a new `Json` object and
// a `bool` identifying success or failure
//
// useful for chained operations when success is important:
//    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
//        log.Println(data)
//    }
func (j *Json) CheckGet(key string) (*Json, bool) {
	j.l.Lock()
	defer j.l.Unlock()
	jj, b := j.sj.CheckGet(key)
	if b {
		return &Json{new(sync.Mutex), jj}, true
	}
	return nil, false
}

// Map type asserts to `map`
func (j *Json) Map() (map[string]interface{}, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Map()
}

// Array type asserts to an `array`
func (j *Json) Array() ([]interface{}, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Array()
}

// Bool type asserts to `bool`
func (j *Json) Bool() (bool, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Bool()
}

// String type asserts to `string`
func (j *Json) String() (string, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.String()
}

// Bytes type asserts to `[]byte`
func (j *Json) Bytes() ([]byte, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Bytes()
}

// StringArray type asserts to an `array` of `string`
func (j *Json) StringArray() ([]string, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.StringArray()
}

// MustArray guarantees the return of a `[]interface{}` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, v := range js.Get("results").MustArray() {
//			fmt.Println(i, v)
//		}
func (j *Json) MustArray(args ...[]interface{}) []interface{} {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustArray()
}

// MustMap guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").MustMap() {
//			fmt.Println(k, v)
//		}
func (j *Json) MustMap(args ...map[string]interface{}) map[string]interface{} {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustMap(args...)
}

// MustString guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").MustString(), js.Get("optional_param").MustString("my_default"))
func (j *Json) MustString(args ...string) string {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustString(args...)
}

// MustStringArray guarantees the return of a `[]string` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, s := range js.Get("results").MustStringArray() {
//			fmt.Println(i, s)
//		}
func (j *Json) MustStringArray(args ...[]string) []string {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustStringArray(args...)
}

// MustInt guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.Get("param1").MustInt(), js.Get("optional_param").MustInt(5150))
func (j *Json) MustInt(args ...int) int {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustInt(args...)
}

// MustFloat64 guarantees the return of a `float64` (with optional default)
//
// useful when you explicitly want a `float64` in a single value return context:
//     myFunc(js.Get("param1").MustFloat64(), js.Get("optional_param").MustFloat64(5.150))
func (j *Json) MustFloat64(args ...float64) float64 {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustFloat64(args...)
}

// MustBool guarantees the return of a `bool` (with optional default)
//
// useful when you explicitly want a `bool` in a single value return context:
//     myFunc(js.Get("param1").MustBool(), js.Get("optional_param").MustBool(true))
func (j *Json) MustBool(args ...bool) bool {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustBool(args...)
}

// MustInt64 guarantees the return of an `int64` (with optional default)
//
// useful when you explicitly want an `int64` in a single value return context:
//     myFunc(js.Get("param1").MustInt64(), js.Get("optional_param").MustInt64(5150))
func (j *Json) MustInt64(args ...int64) int64 {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustInt64(args...)
}

// MustUInt64 guarantees the return of an `uint64` (with optional default)
//
// useful when you explicitly want an `uint64` in a single value return context:
//     myFunc(js.Get("param1").MustUint64(), js.Get("optional_param").MustUint64(5150))
func (j *Json) MustUint64(args ...uint64) uint64 {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.MustUint64(args...)
}

// Float64 coerces into a float64
func (j *Json) Float64() (float64, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Float64()
}

// Int coerces into an int
func (j *Json) Int() (int, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Int()
}

// Int64 coerces into an int64
func (j *Json) Int64() (int64, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Int64()
}

// Uint64 coerces into an uint64
func (j *Json) Uint64() (uint64, error) {
	j.l.Lock()
	defer j.l.Unlock()
	return j.sj.Uint64()
}
