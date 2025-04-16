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
