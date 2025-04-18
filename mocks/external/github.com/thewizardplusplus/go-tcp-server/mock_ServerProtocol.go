// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package tcpServerExternalMocks

import (
	mock "github.com/stretchr/testify/mock"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
)

// NewMockServerProtocol creates a new instance of MockServerProtocol. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockServerProtocol[Req tcpServer.Request, Resp tcpServer.Response](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockServerProtocol[Req, Resp] {
	mock := &MockServerProtocol[Req, Resp]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockServerProtocol is an autogenerated mock type for the ServerProtocol type
type MockServerProtocol[Req tcpServer.Request, Resp tcpServer.Response] struct {
	mock.Mock
}

type MockServerProtocol_Expecter[Req tcpServer.Request, Resp tcpServer.Response] struct {
	mock *mock.Mock
}

func (_m *MockServerProtocol[Req, Resp]) EXPECT() *MockServerProtocol_Expecter[Req, Resp] {
	return &MockServerProtocol_Expecter[Req, Resp]{mock: &_m.Mock}
}

// ExtractToken provides a mock function for the type MockServerProtocol
func (_mock *MockServerProtocol[Req, Resp]) ExtractToken(data []byte, isLatestData bool) (int, []byte, error) {
	ret := _mock.Called(data, isLatestData)

	if len(ret) == 0 {
		panic("no return value specified for ExtractToken")
	}

	var r0 int
	var r1 []byte
	var r2 error
	if returnFunc, ok := ret.Get(0).(func([]byte, bool) (int, []byte, error)); ok {
		return returnFunc(data, isLatestData)
	}
	if returnFunc, ok := ret.Get(0).(func([]byte, bool) int); ok {
		r0 = returnFunc(data, isLatestData)
	} else {
		r0 = ret.Get(0).(int)
	}
	if returnFunc, ok := ret.Get(1).(func([]byte, bool) []byte); ok {
		r1 = returnFunc(data, isLatestData)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}
	if returnFunc, ok := ret.Get(2).(func([]byte, bool) error); ok {
		r2 = returnFunc(data, isLatestData)
	} else {
		r2 = ret.Error(2)
	}
	return r0, r1, r2
}

// MockServerProtocol_ExtractToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExtractToken'
type MockServerProtocol_ExtractToken_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// ExtractToken is a helper method to define mock.On call
//   - data
//   - isLatestData
func (_e *MockServerProtocol_Expecter[Req, Resp]) ExtractToken(data interface{}, isLatestData interface{}) *MockServerProtocol_ExtractToken_Call[Req, Resp] {
	return &MockServerProtocol_ExtractToken_Call[Req, Resp]{Call: _e.mock.On("ExtractToken", data, isLatestData)}
}

func (_c *MockServerProtocol_ExtractToken_Call[Req, Resp]) Run(run func(data []byte, isLatestData bool)) *MockServerProtocol_ExtractToken_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(bool))
	})
	return _c
}

func (_c *MockServerProtocol_ExtractToken_Call[Req, Resp]) Return(offsetToNextToken int, token []byte, err error) *MockServerProtocol_ExtractToken_Call[Req, Resp] {
	_c.Call.Return(offsetToNextToken, token, err)
	return _c
}

func (_c *MockServerProtocol_ExtractToken_Call[Req, Resp]) RunAndReturn(run func(data []byte, isLatestData bool) (int, []byte, error)) *MockServerProtocol_ExtractToken_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// InitialScannerBufferSize provides a mock function for the type MockServerProtocol
func (_mock *MockServerProtocol[Req, Resp]) InitialScannerBufferSize() int {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for InitialScannerBufferSize")
	}

	var r0 int
	if returnFunc, ok := ret.Get(0).(func() int); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(int)
	}
	return r0
}

// MockServerProtocol_InitialScannerBufferSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitialScannerBufferSize'
type MockServerProtocol_InitialScannerBufferSize_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// InitialScannerBufferSize is a helper method to define mock.On call
func (_e *MockServerProtocol_Expecter[Req, Resp]) InitialScannerBufferSize() *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	return &MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp]{Call: _e.mock.On("InitialScannerBufferSize")}
}

func (_c *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp]) Run(run func()) *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp]) Return(n int) *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	_c.Call.Return(n)
	return _c
}

func (_c *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp]) RunAndReturn(run func() int) *MockServerProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// MarshalResponse provides a mock function for the type MockServerProtocol
func (_mock *MockServerProtocol[Req, Resp]) MarshalResponse(response Resp) ([]byte, error) {
	ret := _mock.Called(response)

	if len(ret) == 0 {
		panic("no return value specified for MarshalResponse")
	}

	var r0 []byte
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(Resp) ([]byte, error)); ok {
		return returnFunc(response)
	}
	if returnFunc, ok := ret.Get(0).(func(Resp) []byte); ok {
		r0 = returnFunc(response)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(Resp) error); ok {
		r1 = returnFunc(response)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockServerProtocol_MarshalResponse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalResponse'
type MockServerProtocol_MarshalResponse_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// MarshalResponse is a helper method to define mock.On call
//   - response
func (_e *MockServerProtocol_Expecter[Req, Resp]) MarshalResponse(response interface{}) *MockServerProtocol_MarshalResponse_Call[Req, Resp] {
	return &MockServerProtocol_MarshalResponse_Call[Req, Resp]{Call: _e.mock.On("MarshalResponse", response)}
}

func (_c *MockServerProtocol_MarshalResponse_Call[Req, Resp]) Run(run func(response Resp)) *MockServerProtocol_MarshalResponse_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Resp))
	})
	return _c
}

func (_c *MockServerProtocol_MarshalResponse_Call[Req, Resp]) Return(bytes []byte, err error) *MockServerProtocol_MarshalResponse_Call[Req, Resp] {
	_c.Call.Return(bytes, err)
	return _c
}

func (_c *MockServerProtocol_MarshalResponse_Call[Req, Resp]) RunAndReturn(run func(response Resp) ([]byte, error)) *MockServerProtocol_MarshalResponse_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// MaxTokenSize provides a mock function for the type MockServerProtocol
func (_mock *MockServerProtocol[Req, Resp]) MaxTokenSize() int {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for MaxTokenSize")
	}

	var r0 int
	if returnFunc, ok := ret.Get(0).(func() int); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(int)
	}
	return r0
}

// MockServerProtocol_MaxTokenSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MaxTokenSize'
type MockServerProtocol_MaxTokenSize_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// MaxTokenSize is a helper method to define mock.On call
func (_e *MockServerProtocol_Expecter[Req, Resp]) MaxTokenSize() *MockServerProtocol_MaxTokenSize_Call[Req, Resp] {
	return &MockServerProtocol_MaxTokenSize_Call[Req, Resp]{Call: _e.mock.On("MaxTokenSize")}
}

func (_c *MockServerProtocol_MaxTokenSize_Call[Req, Resp]) Run(run func()) *MockServerProtocol_MaxTokenSize_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockServerProtocol_MaxTokenSize_Call[Req, Resp]) Return(n int) *MockServerProtocol_MaxTokenSize_Call[Req, Resp] {
	_c.Call.Return(n)
	return _c
}

func (_c *MockServerProtocol_MaxTokenSize_Call[Req, Resp]) RunAndReturn(run func() int) *MockServerProtocol_MaxTokenSize_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// ParseRequest provides a mock function for the type MockServerProtocol
func (_mock *MockServerProtocol[Req, Resp]) ParseRequest(token []byte) (Req, error) {
	ret := _mock.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for ParseRequest")
	}

	var r0 Req
	var r1 error
	if returnFunc, ok := ret.Get(0).(func([]byte) (Req, error)); ok {
		return returnFunc(token)
	}
	if returnFunc, ok := ret.Get(0).(func([]byte) Req); ok {
		r0 = returnFunc(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Req)
		}
	}
	if returnFunc, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = returnFunc(token)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockServerProtocol_ParseRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseRequest'
type MockServerProtocol_ParseRequest_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// ParseRequest is a helper method to define mock.On call
//   - token
func (_e *MockServerProtocol_Expecter[Req, Resp]) ParseRequest(token interface{}) *MockServerProtocol_ParseRequest_Call[Req, Resp] {
	return &MockServerProtocol_ParseRequest_Call[Req, Resp]{Call: _e.mock.On("ParseRequest", token)}
}

func (_c *MockServerProtocol_ParseRequest_Call[Req, Resp]) Run(run func(token []byte)) *MockServerProtocol_ParseRequest_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockServerProtocol_ParseRequest_Call[Req, Resp]) Return(v Req, err error) *MockServerProtocol_ParseRequest_Call[Req, Resp] {
	_c.Call.Return(v, err)
	return _c
}

func (_c *MockServerProtocol_ParseRequest_Call[Req, Resp]) RunAndReturn(run func(token []byte) (Req, error)) *MockServerProtocol_ParseRequest_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}
