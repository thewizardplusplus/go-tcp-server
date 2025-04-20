package defaultProtocolModels

import (
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

func TestRequestBuilder_Build(test *testing.T) {
	for _, data := range []struct {
		name    string
		builder *RequestBuilder
		want    Request
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/all parameters",
			builder: NewRequestBuilder().
				SetAction(func() defaultProtocolModelValueTypes.Action {
					value, err := defaultProtocolModelValueTypes.NewAction([]byte("action"))
					require.NoError(test, err)

					return value
				}()).
				SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
					},
				)).
				SetBody(defaultProtocolModelValueTypes.NewBody([]byte("body"))),
			want: Request{
				action: func() defaultProtocolModelValueTypes.Action {
					value, err := defaultProtocolModelValueTypes.NewAction([]byte("action"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
					},
				)),
				body: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("body"))),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/partially empty parameters",
			builder: NewRequestBuilder().
				SetAction(func() defaultProtocolModelValueTypes.Action {
					value, err := defaultProtocolModelValueTypes.NewAction([]byte("action"))
					require.NoError(test, err)

					return value
				}()).
				SetHeaders(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{}, //nolint:lll
				)).
				SetBody(defaultProtocolModelValueTypes.NewBody([]byte{})),
			want: Request{
				action: func() defaultProtocolModelValueTypes.Action {
					value, err := defaultProtocolModelValueTypes.NewAction([]byte("action"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
				body:    mo.None[defaultProtocolModelValueTypes.Body](),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only",
			builder: NewRequestBuilder().
				SetAction(func() defaultProtocolModelValueTypes.Action {
					value, err := defaultProtocolModelValueTypes.NewAction([]byte("action"))
					require.NoError(test, err)

					return value
				}()),
			want: Request{
				action: func() defaultProtocolModelValueTypes.Action {
					value, err := defaultProtocolModelValueTypes.NewAction([]byte("action"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
				body:    mo.None[defaultProtocolModelValueTypes.Body](),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error/without parameters",
			builder: NewRequestBuilder(),
			want:    Request{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := data.builder.Build()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
