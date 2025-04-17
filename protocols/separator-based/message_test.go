package separatorBasedTCPServerProtocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(test *testing.T) {
	type args struct {
		data   []byte
		params SeparationParams
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Message
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/minimal message",
			args: args{
				data: []byte("dummy||"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want: Message{
				Introduction: []byte("dummy"),
				Headers:      make(map[string][]byte),
				Body:         []byte{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/full message",
			args: args{
				data: []byte("introduction|one=two&three=four|body"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want: Message{
				Introduction: []byte("introduction"),
				Headers: map[string][]byte{
					"6f6e65":     []byte("two"),
					"7468726565": []byte("four"),
				},
				Body: []byte("body"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped data/regular data",
			args: args{
				data: []byte("intro%64%75ction|o%6ee=t%77o&t%68%72%65e=f%6f%75r|b%6f%64y"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want: Message{
				Introduction: []byte("introduction"),
				Headers: map[string][]byte{
					"6f6e65":     []byte("two"),
					"7468726565": []byte("four"),
				},
				Body: []byte("body"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped data/separators",
			args: args{
				data: []byte(
					"dummy%7cintroduction%0a|" +
						"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
						"dummy%7cbody%0a",
				),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want: Message{
				Introduction: []byte("dummy|introduction\n"),
				Headers: map[string][]byte{
					"64756d6d793d6f6e650a":     []byte("dummy&two\n"),
					"64756d6d793d74687265650a": []byte("dummy&four\n"),
				},
				Body: []byte("dummy|body\n"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/invalid message part count",
			args: args{
				data: []byte("dummy"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/introduction cannot be empty",
			args: args{
				data: []byte("|dummy #1|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the introduction",
			args: args{
				data: []byte("dummy %xx|dummy #1|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/header has no key-value separator",
			args: args{
				data: []byte("dummy #0|dummy #1|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/header key cannot be empty",
			args: args{
				data: []byte("dummy #0|=value|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the header key",
			args: args{
				data: []byte("dummy #0|key %xx=value|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/header value cannot be empty",
			args: args{
				data: []byte("dummy #0|key=|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the header value",
			args: args{
				data: []byte("dummy #0|key=value %xx|dummy #2"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the body",
			args: args{
				data: []byte("dummy #0|key=value|dummy %xx"),
				params: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want:    Message{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := ParseMessage(data.args.data, data.args.params)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
