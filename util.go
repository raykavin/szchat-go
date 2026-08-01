package szchat

import "encoding/json"

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
