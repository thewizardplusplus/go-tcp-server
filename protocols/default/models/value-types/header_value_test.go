package defaultProtocolModelValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeaderValue(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name    string
		args    args
		want    HeaderValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: HeaderValue{
				rawValue: []byte("dummy"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/empty",
			args: args{
				rawValue: []byte{},
			},
			want:    HeaderValue{},
			wantErr: assert.Error,
		},
		{
			name: "error/nil",
			args: args{
				rawValue: nil,
			},
			want:    HeaderValue{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewHeaderValue(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestMustNewHeaderValue(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name      string
		args      args
		want      HeaderValue
		wantPanic assert.PanicAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: HeaderValue{
				rawValue: []byte("dummy"),
			},
			wantPanic: assert.NotPanics,
		},
		{
			name: "error",
			args: args{
				rawValue: []byte{},
			},
			want:      HeaderValue{},
			wantPanic: assert.Panics,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var got HeaderValue
			data.wantPanic(test, func() {
				got = MustNewHeaderValue(data.args.rawValue)
			})

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHeaderValue_ToBytes(test *testing.T) {
	type fields struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "success",
			fields: fields{
				rawValue: []byte("dummy"),
			},
			want: []byte("dummy"),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HeaderValue{
				rawValue: data.fields.rawValue,
			}
			got := value.ToBytes()

			assert.Equal(test, data.want, got)
		})
	}
}
