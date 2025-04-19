package defaultProtocolModelValueTypes

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderKey_comparable(test *testing.T) {
	assert.True(test, reflect.ValueOf(HeaderKey{}).Comparable())
}

func TestNewHeaderKey(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name    string
		args    args
		want    HeaderKey
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: HeaderKey{
				encodedRawValue: "64756d6d79",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/empty",
			args: args{
				rawValue: []byte{},
			},
			want:    HeaderKey{},
			wantErr: assert.Error,
		},
		{
			name: "error/nil",
			args: args{
				rawValue: nil,
			},
			want:    HeaderKey{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewHeaderKey(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestMustNewHeaderKey(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name      string
		args      args
		want      HeaderKey
		wantPanic assert.PanicAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: HeaderKey{
				encodedRawValue: "64756d6d79",
			},
			wantPanic: assert.NotPanics,
		},
		{
			name: "error",
			args: args{
				rawValue: []byte{},
			},
			want:      HeaderKey{},
			wantPanic: assert.Panics,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var got HeaderKey
			data.wantPanic(test, func() {
				got = MustNewHeaderKey(data.args.rawValue)
			})

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHeaderKey_ToBytes(test *testing.T) {
	type fields struct {
		encodedRawValue string
	}

	for _, data := range []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				encodedRawValue: "64756d6d79",
			},
			want:    []byte("dummy"),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				encodedRawValue: "invalid",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HeaderKey{
				encodedRawValue: data.fields.encodedRawValue,
			}
			got, err := value.ToBytes()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
