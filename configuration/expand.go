package configuration

func Expand(s string, mapping func(idx int, key string) string) string {
	s, _ = expand(s, func(idx int, key string) (string, error) {
		return mapping(idx, key), nil
	})
	return s
}

func ExpandError(s string, mapping func(idx int, key string) (string, error)) (string, error) {
	return expand(s, mapping)
}

func ExpandBytes(data []byte, mapping func(idx int, key string) []byte) []byte {
	data, _ = expandBytes(data, func(idx int, key string) ([]byte, error) {
		return mapping(idx, key), nil
	})
	return data
}

func ExpandBytesError(data []byte, mapping func(idx int, key string) ([]byte, error)) ([]byte, error) {
	return expandBytes(data, mapping)
}

func expand(s string, mapping func(idx int, key string) (string, error)) (string, error) {
	var buf = make([]byte, 0, 2*len(s))
	var i = 0
	var idx = 0
	var left = -1
	for j := 0; j < len(s); j++ {
		if s[j] == '{' {
			left = j
			continue
		}
		if s[j] != '}' {
			continue
		}
		if left < 0 {
			continue
		}
		key := s[left+1 : j]
		buf = append(buf, s[i:left]...)
		val, err := mapping(idx, key)
		if err != nil {
			return "", err
		}
		buf = append(buf, val...)
		i = j + 1
		left = -1
		idx++
	}
	buf = append(buf, s[i:]...)
	return string(buf), nil
}

func expandBytes(s []byte, mapping func(idx int, key string) ([]byte, error)) ([]byte, error) {
	var buf = make([]byte, 0, 2*len(s))
	var i = 0
	var idx = 0
	var left = -1
	for j := 0; j < len(s); j++ {
		if s[j] == '{' {
			left = j
			continue
		}
		if s[j] != '}' {
			continue
		}
		if left < 0 {
			continue
		}
		key := s[left+1 : j]
		buf = append(buf, s[i:left]...)
		val, err := mapping(idx, string(key))
		if err != nil {
			return nil, err
		}
		buf = append(buf, val...)
		i = j + 1
		left = -1
		idx++
	}
	buf = append(buf, s[i:]...)
	return buf, nil
}

func SliceGetter(data ...string) func(int, string) string {
	return func(idx int, _ string) string {
		if idx >= len(data) {
			return ""
		}
		m := data[idx]
		return m
	}
}

func MapGetter(m map[string]string) func(int, string) string {
	return func(_ int, key string) string {
		return m[key]
	}
}
