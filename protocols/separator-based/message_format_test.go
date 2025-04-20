package separatorBasedProtocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	defaultProtocol "github.com/thewizardplusplus/go-tcp-server/protocols/default"
	defaultProtocolModels "github.com/thewizardplusplus/go-tcp-server/protocols/default/models"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

func TestMessageFormat_interface(test *testing.T) {
	assert.Implements(test, (*defaultProtocol.MessageFormat)(nil), MessageFormat{})
}

func TestNewMessageFormat(test *testing.T) {
	type args struct {
		options SeparationParams
	}

	for _, data := range []struct {
		name string
		args args
		want MessageFormat
	}{
		{
			name: "success",
			args: args{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want: MessageFormat{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewMessageFormat(data.args.options)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestMessageFormat_ParseMessage(test *testing.T) {
	type fields struct {
		options SeparationParams
	}
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    defaultProtocolModels.Message
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/minimal message",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy||"),
			},
			want: func() defaultProtocolModels.Message {
				introduction, err :=
					defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
				require.NoError(test, err)

				message, err := defaultProtocolModels.NewMessageBuilder().
					SetIntroduction(introduction).
					Build()
				require.NoError(test, err)

				return message
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success/full message",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("introduction|one=two&three=four|body"),
			},
			want: func() defaultProtocolModels.Message {
				introduction, err :=
					defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
				require.NoError(test, err)

				message, err := defaultProtocolModels.NewMessageBuilder().
					SetIntroduction(introduction).
					SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
						map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
						},
					)).
					SetBody(defaultProtocolModelValueTypes.NewBody([]byte("body"))).
					Build()
				require.NoError(test, err)

				return message
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped data/regular data",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("intro%64%75ction|o%6ee=t%77o&t%68%72%65e=f%6f%75r|b%6f%64y"),
			},
			want: func() defaultProtocolModels.Message {
				introduction, err :=
					defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
				require.NoError(test, err)

				message, err := defaultProtocolModels.NewMessageBuilder().
					SetIntroduction(introduction).
					SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
						map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
						},
					)).
					SetBody(defaultProtocolModelValueTypes.NewBody([]byte("body"))).
					Build()
				require.NoError(test, err)

				return message
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped data/separators",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte(
					"dummy%7cintroduction%0a|" +
						"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
						"dummy%7cbody%0a",
				),
			},
			want: func() defaultProtocolModels.Message {
				introduction, err :=
					defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy|introduction\n")) //nolint:lll
				require.NoError(test, err)

				message, err := defaultProtocolModels.NewMessageBuilder().
					SetIntroduction(introduction).
					SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
						map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("dummy=one\n")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("dummy&two\n")),  //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("dummy=three\n")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("dummy&four\n")), //nolint:lll
						},
					)).
					SetBody(defaultProtocolModelValueTypes.NewBody([]byte("dummy|body\n"))).
					Build()
				require.NoError(test, err)

				return message
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error/invalid message part count",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/introduction cannot be empty",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("|dummy #1|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the introduction",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy %xx|dummy #1|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/header has no key-value separator",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy #0|dummy #1|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/header key cannot be empty",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy #0|=value|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the header key",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy #0|key %xx=value|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/header value cannot be empty",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy #0|key=|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the header value",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy #0|key=value %xx|dummy #2"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to unescape separators in the body",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				data: []byte("dummy #0|key=value|dummy %xx"),
			},
			want:    defaultProtocolModels.Message{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			format := MessageFormat{
				options: data.fields.options,
			}
			got, err := format.ParseMessage(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestMessageFormat_MarshalMessage(test *testing.T) {
	type fields struct {
		options SeparationParams
	}
	type args struct {
		message defaultProtocolModels.Message
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/minimal message",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				message: func() defaultProtocolModels.Message {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					return message
				}(),
			},
			want:    []byte("dummy||"),
			wantErr: assert.NoError,
		},
		{
			name: "success/full message",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				message: func() defaultProtocolModels.Message {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
							map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
								defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
								defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
							},
						)).
						SetBody(defaultProtocolModelValueTypes.NewBody([]byte("body"))).
						Build()
					require.NoError(test, err)

					return message
				}(),
			},
			want:    []byte("introduction|one=two&three=four|body"),
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped separators",
			fields: fields{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			args: args{
				message: func() defaultProtocolModels.Message {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy|introduction\n")) //nolint:lll
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
							map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
								defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("dummy=one\n")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("dummy&two\n")),  //nolint:lll
								defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("dummy=three\n")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("dummy&four\n")), //nolint:lll
							},
						)).
						SetBody(defaultProtocolModelValueTypes.NewBody([]byte("dummy|body\n"))).
						Build()
					require.NoError(test, err)

					return message
				}(),
			},
			want: []byte(
				"dummy%7cintroduction%0a|" +
					"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
					"dummy%7cbody%0a",
			),
			wantErr: assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			format := MessageFormat{
				options: data.fields.options,
			}
			got, err := format.MarshalMessage(data.args.message)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
