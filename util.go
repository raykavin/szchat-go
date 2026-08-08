package szchat

import (
	"encoding/json"
	"reflect"
	"strings"
)

// StringSlice decodes a JSON array of strings, a single JSON
// string, or null into a []string: some SZChat fields that can hold
// multiple values are only wrapped in an array once there's more than one,
// coming back as a bare string otherwise.
type StringSlice []string

func (s *StringSlice) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if len(data) > 0 && data[0] == '[' {
		var slice []string
		if err := json.Unmarshal(data, &slice); err != nil {
			return err
		}
		*s = StringSlice(slice)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		*s = nil
		return nil
	}
	*s = StringSlice{str}
	return nil
}

// marshalWithExtra marshals v and shallow-merges extra on top of the result,
// letting request structs expose a fixed set of well-known fields while
// still accepting arbitrary/tenant-specific keys (e.g. custom contact
// channels) without needing a dedicated field per key.
func marshalWithExtra(v any, extra map[string]any) ([]byte, error) {
	base, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(extra) == 0 {
		return base, nil
	}

	merged := map[string]any{}
	if err := json.Unmarshal(base, &merged); err != nil {
		return nil, err
	}
	for k, v := range extra {
		merged[k] = v
	}
	return json.Marshal(merged)
}

// unmarshalWithExtra decodes data into dst using the default struct decode,
// then returns every top-level JSON key in data that dst's own json tags
// don't account for. It is marshalWithExtra's counterpart: known fields in,
// unknown/dynamic fields out, so a struct's UnmarshalJSON can keep tenant-
// specific custom keys instead of silently dropping them.
func unmarshalWithExtra(data []byte, dst any) (map[string]json.RawMessage, error) {
	if err := json.Unmarshal(data, dst); err != nil {
		return nil, err
	}

	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, err
	}

	for _, name := range jsonFieldNames(dst) {
		delete(all, name)
	}

	return all, nil
}

// jsonFieldNames returns the JSON object keys that v's struct tags map to,
// read via reflection so the known/unknown split in unmarshalWithExtra
// can't drift out of sync as fields are added to v's type over time.
func jsonFieldNames(v any) []string {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	names := make([]string, 0, t.NumField())
	for i := range t.NumField() {
		tag, ok := t.Field(i).Tag.Lookup("json")
		if !ok {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "" || name == "-" {
			continue
		}
		names = append(names, name)
	}
	return names
}
