package defaultProtocolModels

import (
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

func TestMessage_Introduction(test *testing.T) {
	type fields struct {
		introduction defaultProtocolModelValueTypes.Introduction
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   defaultProtocolModelValueTypes.Introduction
	}{
		{
			name: "success",
			fields: fields{
				introduction: func() defaultProtocolModelValueTypes.Introduction {
					value, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
			},
			want: func() defaultProtocolModelValueTypes.Introduction {
				value, err :=
					defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
				require.NoError(test, err)

				return value
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			model := Message{
				introduction: data.fields.introduction,
			}
			got := model.Introduction()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestMessage_Headers(test *testing.T) {
	type fields struct {
		headers mo.Option[defaultProtocolModelValueTypes.Headers]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   mo.Option[defaultProtocolModelValueTypes.Headers]
	}{
		{
			name: "success/is present",
			fields: fields{
				headers: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
					},
				)),
			},
			want: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
				map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
					defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
					defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
				},
			)),
		},
		{
			name: "success/is absent",
			fields: fields{
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
			},
			want: mo.None[defaultProtocolModelValueTypes.Headers](),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			model := Message{
				headers: data.fields.headers,
			}
			got := model.Headers()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestMessage_Body(test *testing.T) {
	type fields struct {
		body mo.Option[defaultProtocolModelValueTypes.Body]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   mo.Option[defaultProtocolModelValueTypes.Body]
	}{
		{
			name: "success/is present",
			fields: fields{
				body: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("dummy"))),
			},
			want: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("dummy"))),
		},
		{
			name: "success/is absent",
			fields: fields{
				body: mo.None[defaultProtocolModelValueTypes.Body](),
			},
			want: mo.None[defaultProtocolModelValueTypes.Body](),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			model := Message{
				body: data.fields.body,
			}
			got := model.Body()

			assert.Equal(test, data.want, got)
		})
	}
}
