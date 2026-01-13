package configuration

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpand(t *testing.T) {
	m := map[string]string{"FOO": "BAR"}
	assert.Equal(t, "BAR", Expand("{FOO}", MapGetter(m)))
	assert.Equal(t, "{BAR}", Expand("{{FOO}}", MapGetter(m)))
	assert.Equal(t, "", Expand("{FOOO}", MapGetter(m)))
	assert.Equal(t, "", Expand("{}", MapGetter(m)))
	assert.Equal(t, "{", Expand("{{}", MapGetter(m)))
	assert.Equal(t, "}", Expand("{}}", MapGetter(m)))
	assert.Equal(t, "}{", Expand("{}}{{}", MapGetter(m)))
	assert.Equal(t, "{}", Expand("{{}}", MapGetter(m)))
	assert.Equal(t, "{{}}", Expand("{{{}}}", MapGetter(m)))
	assert.Equal(t, "$BAR", Expand("${FOO}", MapGetter(m)))
}

func BenchmarkExpand(b *testing.B) {
	m := map[string]string{"FOO": "BAR"}
	data := strings.Repeat("A", 100) + strings.Repeat("${FOO}", 50) + strings.Repeat("A", 100)
	b.ResetTimer()
	b.Run("expand", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Expand(data, MapGetter(m))
		}
	})
	b.Run("std", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			os.Expand(data, func(s string) string { return m[s] })
		}
	})
}
