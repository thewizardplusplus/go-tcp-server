package defaultProtocolModelValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntroduction(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Introduction
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: Introduction{
				rawValue: []byte("dummy"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/empty",
			args: args{
				rawValue: []byte{},
			},
			want:    Introduction{},
			wantErr: assert.Error,
		},
		{
			name: "error/nil",
			args: args{
				rawValue: nil,
			},
			want:    Introduction{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewIntroduction(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestIntroduction_ToBytes(test *testing.T) {
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
			value := Introduction{
				rawValue: data.fields.rawValue,
			}
			got := value.ToBytes()

			assert.Equal(test, data.want, got)
		})
	}
}
