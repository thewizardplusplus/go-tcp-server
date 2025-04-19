package defaultProtocolModelValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBody(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name string
		args args
		want Body
	}{
		{
			name: "success/non-empty",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: Body{
				rawValue: []byte("dummy"),
			},
		},
		{
			name: "success/empty",
			args: args{
				rawValue: []byte{},
			},
			want: Body{
				rawValue: []byte{},
			},
		},
		{
			name: "success/nil",
			args: args{
				rawValue: nil,
			},
			want: Body{
				rawValue: []byte{},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewBody(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestBody_ToBytes(test *testing.T) {
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
			value := Body{
				rawValue: data.fields.rawValue,
			}
			got := value.ToBytes()

			assert.Equal(test, data.want, got)
		})
	}
}
