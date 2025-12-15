package mesh

type Metadata map[string][]string

func (m Metadata) Join(md Metadata) {
	for k, v := range md {
		m[k] = append(m[k], v...)
	}
}

func (m Metadata) Get(k string) string {
	if len(m[k]) == 0 {
		return ""
	}
	return m[k][0]
}
