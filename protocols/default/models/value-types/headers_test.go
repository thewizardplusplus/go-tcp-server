package defaultProtocolModelValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeaders(test *testing.T) {
	type args struct {
		rawValue map[HeaderKey]HeaderValue
	}

	for _, data := range []struct {
		name string
		args args
		want Headers
	}{
		{
			name: "success/non-empty",
			args: args{
				rawValue: map[HeaderKey]HeaderValue{
					MustNewHeaderKey([]byte("one")):   MustNewHeaderValue([]byte("two")),
					MustNewHeaderKey([]byte("three")): MustNewHeaderValue([]byte("four")),
				},
			},
			want: Headers{
				rawValue: map[HeaderKey]HeaderValue{
					MustNewHeaderKey([]byte("one")):   MustNewHeaderValue([]byte("two")),
					MustNewHeaderKey([]byte("three")): MustNewHeaderValue([]byte("four")),
				},
			},
		},
		{
			name: "success/empty",
			args: args{
				rawValue: map[HeaderKey]HeaderValue{},
			},
			want: Headers{
				rawValue: map[HeaderKey]HeaderValue{},
			},
		},
		{
			name: "success/nil",
			args: args{
				rawValue: nil,
			},
			want: Headers{
				rawValue: map[HeaderKey]HeaderValue{},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewHeaders(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHeaders_ToMap(test *testing.T) {
	type fields struct {
		rawValue map[HeaderKey]HeaderValue
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   map[HeaderKey]HeaderValue
	}{
		{
			name: "success",
			fields: fields{
				rawValue: map[HeaderKey]HeaderValue{
					MustNewHeaderKey([]byte("one")):   MustNewHeaderValue([]byte("two")),
					MustNewHeaderKey([]byte("three")): MustNewHeaderValue([]byte("four")),
				},
			},
			want: map[HeaderKey]HeaderValue{
				MustNewHeaderKey([]byte("one")):   MustNewHeaderValue([]byte("two")),
				MustNewHeaderKey([]byte("three")): MustNewHeaderValue([]byte("four")),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Headers{
				rawValue: data.fields.rawValue,
			}
			got := value.ToMap()

			assert.Equal(test, data.want, got)
		})
	}
}
