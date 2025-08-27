package url

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateFunctionPrepareURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should return url with scheme https if got url without scheme",
			args: args{url: "example.com/path/?param=value&key=value"},
			want: "https://example.com/path/?param=value&key=value",
		},
		{
			name: "Should return url without canges if got url with scheme",
			args: args{url: "https://www.yandex.ru/path/?param=value"},
			want: "https://www.yandex.ru/path/?param=value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, prepareURL(tt.args.url))
		})
	}
}
