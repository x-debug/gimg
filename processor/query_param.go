package processor

import (
	"strconv"
	"strings"
)

const (
	JoinSplit     = "&"
	KeyValueSplit = "="
)

type HttpFinger interface {
	ActionFinger() string
}

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

//GetInt return integer value or default value
func (ps Params) GetInt(name string, defaultValue int) int {
	va, ok := ps.Get(name)
	if !ok {
		return defaultValue
	}

	intVa, err := strconv.Atoi(va)
	if err != nil {
		return defaultValue
	}

	return intVa
}

func (ps Params) GetFloat64(name string, defaultValue float64) float64 {
	va, ok := ps.Get(name)
	if !ok {
		return defaultValue
	}

	f64, err := strconv.ParseFloat(va, 64)
	if err != nil {
		return defaultValue
	}

	return f64
}

func (ps Params) GetFloat32(name string, defaultValue float32) float32 {
	va, ok := ps.Get(name)
	if !ok {
		return defaultValue
	}

	f32, err := strconv.ParseFloat(va, 32)
	if err != nil {
		return defaultValue
	}

	return float32(f32)
}

//GetString return string value or default value
func (ps Params) GetString(name string, defaultValue string) string {
	va, ok := ps.Get(name)
	if !ok {
		return defaultValue
	}

	return va
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

//NewParams convert url query to params
func NewParams(query string) Params {
	paris := strings.Split(query, JoinSplit)

	params := make(Params, 0)
	for _, pair := range paris {
		keyValues := strings.Split(pair, KeyValueSplit)
		if len(keyValues) > 1 {
			params = append(params, Param{Key: keyValues[0], Value: keyValues[1]})
		} else {
			params = append(params, Param{Key: keyValues[0], Value: ""})
		}
	}

	return params
}
