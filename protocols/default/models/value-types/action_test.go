package defaultProtocolModelValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAction(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Action
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: Action{
				rawValue: []byte("dummy"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/empty",
			args: args{
				rawValue: []byte{},
			},
			want:    Action{},
			wantErr: assert.Error,
		},
		{
			name: "error/nil",
			args: args{
				rawValue: nil,
			},
			want:    Action{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewAction(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestAction_ToBytes(test *testing.T) {
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
			value := Action{
				rawValue: data.fields.rawValue,
			}
			got := value.ToBytes()

			assert.Equal(test, data.want, got)
		})
	}
}
