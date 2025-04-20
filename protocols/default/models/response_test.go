package defaultProtocolModels

import (
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

func TestNewResponseFromMessage(test *testing.T) {
	type args struct {
		message Message
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Response
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/all parameters",
			args: args{
				message: Message{
					introduction: func() defaultProtocolModelValueTypes.Introduction {
						value, err :=
							defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
						require.NoError(test, err)

						return value
					}(),
					headers: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
						map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
							defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
						},
					)),
					body: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("dummy"))),
				},
			},
			want: Response{
				status: func() defaultProtocolModelValueTypes.Status {
					value, err := defaultProtocolModelValueTypes.NewStatus([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
					},
				)),
				body: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("dummy"))),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only",
			args: args{
				message: Message{
					introduction: func() defaultProtocolModelValueTypes.Introduction {
						value, err :=
							defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
						require.NoError(test, err)

						return value
					}(),
					headers: mo.None[defaultProtocolModelValueTypes.Headers](),
					body:    mo.None[defaultProtocolModelValueTypes.Body](),
				},
			},
			want: Response{
				status: func() defaultProtocolModelValueTypes.Status {
					value, err := defaultProtocolModelValueTypes.NewStatus([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
				body:    mo.None[defaultProtocolModelValueTypes.Body](),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				message: Message{
					introduction: defaultProtocolModelValueTypes.Introduction{},
					headers:      mo.None[defaultProtocolModelValueTypes.Headers](),
					body:         mo.None[defaultProtocolModelValueTypes.Body](),
				},
			},
			want:    Response{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewResponseFromMessage(data.args.message)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestResponse_Status(test *testing.T) {
	type fields struct {
		status defaultProtocolModelValueTypes.Status
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   defaultProtocolModelValueTypes.Status
	}{
		{
			name: "success",
			fields: fields{
				status: func() defaultProtocolModelValueTypes.Status {
					value, err := defaultProtocolModelValueTypes.NewStatus([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
			},
			want: func() defaultProtocolModelValueTypes.Status {
				value, err := defaultProtocolModelValueTypes.NewStatus([]byte("dummy"))
				require.NoError(test, err)

				return value
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			model := Response{
				status: data.fields.status,
			}
			got := model.Status()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestResponse_Headers(test *testing.T) {
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
			model := Response{
				headers: data.fields.headers,
			}
			got := model.Headers()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestResponse_Body(test *testing.T) {
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
			model := Response{
				body: data.fields.body,
			}
			got := model.Body()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestResponse_ToMessage(test *testing.T) {
	type fields struct {
		status  defaultProtocolModelValueTypes.Status
		headers mo.Option[defaultProtocolModelValueTypes.Headers]
		body    mo.Option[defaultProtocolModelValueTypes.Body]
	}

	for _, data := range []struct {
		name    string
		fields  fields
		want    Message
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/all parameters",
			fields: fields{
				status: func() defaultProtocolModelValueTypes.Status {
					value, err := defaultProtocolModelValueTypes.NewStatus([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
					},
				)),
				body: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("dummy"))),
			},
			want: Message{
				introduction: func() defaultProtocolModelValueTypes.Introduction {
					value, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.Some(defaultProtocolModelValueTypes.NewHeaders(
					map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue{ //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("one")):   defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("two")),  //nolint:lll
						defaultProtocolModelValueTypes.MustNewHeaderKey([]byte("three")): defaultProtocolModelValueTypes.MustNewHeaderValue([]byte("four")), //nolint:lll
					},
				)),
				body: mo.Some(defaultProtocolModelValueTypes.NewBody([]byte("dummy"))),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only",
			fields: fields{
				status: func() defaultProtocolModelValueTypes.Status {
					value, err := defaultProtocolModelValueTypes.NewStatus([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
				body:    mo.None[defaultProtocolModelValueTypes.Body](),
			},
			want: Message{
				introduction: func() defaultProtocolModelValueTypes.Introduction {
					value, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("dummy"))
					require.NoError(test, err)

					return value
				}(),
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
				body:    mo.None[defaultProtocolModelValueTypes.Body](),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to construct the introduction",
			fields: fields{
				status:  defaultProtocolModelValueTypes.Status{},
				headers: mo.None[defaultProtocolModelValueTypes.Headers](),
				body:    mo.None[defaultProtocolModelValueTypes.Body](),
			},
			want:    Message{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			model := Response{
				status:  data.fields.status,
				headers: data.fields.headers,
				body:    data.fields.body,
			}
			got, err := model.ToMessage()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
