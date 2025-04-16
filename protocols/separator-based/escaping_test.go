package separatorBasedTCPServerProtocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeSeparators(test *testing.T) {
	type args struct {
		data       []byte
		separators [][]byte
	}

	for _, data := range []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "success/without separators",
			args: args{
				data:       []byte("dummy"),
				separators: nil,
			},
			want: []byte("dummy"),
		},
		{
			name: "success/with a single separator",
			args: args{
				data:       []byte("one:=two;three:=four;"),
				separators: [][]byte{[]byte(":=")},
			},
			want: []byte("one%3a%3dtwo;three%3a%3dfour;"),
		},
		{
			name: "success/with several separators",
			args: args{
				data:       []byte("one:=two;three<-four;five:=six;seven<-eight;"),
				separators: [][]byte{[]byte(":="), []byte("<-")},
			},
			want: []byte("one%3a%3dtwo;three%3c%2dfour;five%3a%3dsix;seven%3c%2deight;"),
		},
		{
			name: "success/with the percent character",
			args: args{
				data:       []byte("one:=23%;two:=42%;"),
				separators: [][]byte{[]byte(":=")},
			},
			want: []byte("one%3a%3d23%25;two%3a%3d42%25;"),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := EscapeSeparators(data.args.data, data.args.separators)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestUnescapeSeparators(test *testing.T) {
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/without separators",
			args: args{
				data: []byte("dummy"),
			},
			want:    []byte("dummy"),
			wantErr: assert.NoError,
		},
		{
			name: "success/with a single separator type",
			args: args{
				data: []byte("one%3a%3dtwo;three%3a%3dfour;"),
			},
			want:    []byte("one:=two;three:=four;"),
			wantErr: assert.NoError,
		},
		{
			name: "success/with several separator types",
			args: args{
				data: []byte(
					"one%3a%3dtwo;three%3c%2dfour;five%3a%3dsix;seven%3c%2deight;",
				),
			},
			want:    []byte("one:=two;three<-four;five:=six;seven<-eight;"),
			wantErr: assert.NoError,
		},
		{
			name: "success/with the percent character",
			args: args{
				data: []byte("one%3a%3d23%25;two%3a%3d42%25;"),
			},
			want:    []byte("one:=23%;two:=42%;"),
			wantErr: assert.NoError,
		},
		{
			name: "error/not enough characters to decode/zero characters",
			args: args{
				data: []byte("one%3a%3dtwo;three%3a%"),
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error/not enough characters to decode/single character",
			args: args{
				data: []byte("one%3a%3dtwo;three%3a%3"),
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error/unable to decode the sequence",
			args: args{
				data: []byte("one%3a%xxtwo;three%3a%3dfour;"),
			},
			want:    nil,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := UnescapeSeparators(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
