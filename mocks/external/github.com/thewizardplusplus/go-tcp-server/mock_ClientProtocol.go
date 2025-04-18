// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package tcpServerExternalMocks

import (
	mock "github.com/stretchr/testify/mock"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
)

// NewMockClientProtocol creates a new instance of MockClientProtocol. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClientProtocol[Req tcpServer.Request, Resp tcpServer.Response](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClientProtocol[Req, Resp] {
	mock := &MockClientProtocol[Req, Resp]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockClientProtocol is an autogenerated mock type for the ClientProtocol type
type MockClientProtocol[Req tcpServer.Request, Resp tcpServer.Response] struct {
	mock.Mock
}

type MockClientProtocol_Expecter[Req tcpServer.Request, Resp tcpServer.Response] struct {
	mock *mock.Mock
}

func (_m *MockClientProtocol[Req, Resp]) EXPECT() *MockClientProtocol_Expecter[Req, Resp] {
	return &MockClientProtocol_Expecter[Req, Resp]{mock: &_m.Mock}
}

// ExtractToken provides a mock function for the type MockClientProtocol
func (_mock *MockClientProtocol[Req, Resp]) ExtractToken(data []byte, isLatestData bool) (int, []byte, error) {
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

// MockClientProtocol_ExtractToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExtractToken'
type MockClientProtocol_ExtractToken_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// ExtractToken is a helper method to define mock.On call
//   - data
//   - isLatestData
func (_e *MockClientProtocol_Expecter[Req, Resp]) ExtractToken(data interface{}, isLatestData interface{}) *MockClientProtocol_ExtractToken_Call[Req, Resp] {
	return &MockClientProtocol_ExtractToken_Call[Req, Resp]{Call: _e.mock.On("ExtractToken", data, isLatestData)}
}

func (_c *MockClientProtocol_ExtractToken_Call[Req, Resp]) Run(run func(data []byte, isLatestData bool)) *MockClientProtocol_ExtractToken_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(bool))
	})
	return _c
}

func (_c *MockClientProtocol_ExtractToken_Call[Req, Resp]) Return(offsetToNextToken int, token []byte, err error) *MockClientProtocol_ExtractToken_Call[Req, Resp] {
	_c.Call.Return(offsetToNextToken, token, err)
	return _c
}

func (_c *MockClientProtocol_ExtractToken_Call[Req, Resp]) RunAndReturn(run func(data []byte, isLatestData bool) (int, []byte, error)) *MockClientProtocol_ExtractToken_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// InitialScannerBufferSize provides a mock function for the type MockClientProtocol
func (_mock *MockClientProtocol[Req, Resp]) InitialScannerBufferSize() int {
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

// MockClientProtocol_InitialScannerBufferSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitialScannerBufferSize'
type MockClientProtocol_InitialScannerBufferSize_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// InitialScannerBufferSize is a helper method to define mock.On call
func (_e *MockClientProtocol_Expecter[Req, Resp]) InitialScannerBufferSize() *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	return &MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp]{Call: _e.mock.On("InitialScannerBufferSize")}
}

func (_c *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp]) Run(run func()) *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp]) Return(n int) *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	_c.Call.Return(n)
	return _c
}

func (_c *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp]) RunAndReturn(run func() int) *MockClientProtocol_InitialScannerBufferSize_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// MarshalRequest provides a mock function for the type MockClientProtocol
func (_mock *MockClientProtocol[Req, Resp]) MarshalRequest(request Req) ([]byte, error) {
	ret := _mock.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for MarshalRequest")
	}

	var r0 []byte
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(Req) ([]byte, error)); ok {
		return returnFunc(request)
	}
	if returnFunc, ok := ret.Get(0).(func(Req) []byte); ok {
		r0 = returnFunc(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(Req) error); ok {
		r1 = returnFunc(request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockClientProtocol_MarshalRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalRequest'
type MockClientProtocol_MarshalRequest_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// MarshalRequest is a helper method to define mock.On call
//   - request
func (_e *MockClientProtocol_Expecter[Req, Resp]) MarshalRequest(request interface{}) *MockClientProtocol_MarshalRequest_Call[Req, Resp] {
	return &MockClientProtocol_MarshalRequest_Call[Req, Resp]{Call: _e.mock.On("MarshalRequest", request)}
}

func (_c *MockClientProtocol_MarshalRequest_Call[Req, Resp]) Run(run func(request Req)) *MockClientProtocol_MarshalRequest_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Req))
	})
	return _c
}

func (_c *MockClientProtocol_MarshalRequest_Call[Req, Resp]) Return(bytes []byte, err error) *MockClientProtocol_MarshalRequest_Call[Req, Resp] {
	_c.Call.Return(bytes, err)
	return _c
}

func (_c *MockClientProtocol_MarshalRequest_Call[Req, Resp]) RunAndReturn(run func(request Req) ([]byte, error)) *MockClientProtocol_MarshalRequest_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// MaxTokenSize provides a mock function for the type MockClientProtocol
func (_mock *MockClientProtocol[Req, Resp]) MaxTokenSize() int {
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

// MockClientProtocol_MaxTokenSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MaxTokenSize'
type MockClientProtocol_MaxTokenSize_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// MaxTokenSize is a helper method to define mock.On call
func (_e *MockClientProtocol_Expecter[Req, Resp]) MaxTokenSize() *MockClientProtocol_MaxTokenSize_Call[Req, Resp] {
	return &MockClientProtocol_MaxTokenSize_Call[Req, Resp]{Call: _e.mock.On("MaxTokenSize")}
}

func (_c *MockClientProtocol_MaxTokenSize_Call[Req, Resp]) Run(run func()) *MockClientProtocol_MaxTokenSize_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockClientProtocol_MaxTokenSize_Call[Req, Resp]) Return(n int) *MockClientProtocol_MaxTokenSize_Call[Req, Resp] {
	_c.Call.Return(n)
	return _c
}

func (_c *MockClientProtocol_MaxTokenSize_Call[Req, Resp]) RunAndReturn(run func() int) *MockClientProtocol_MaxTokenSize_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}

// ParseResponse provides a mock function for the type MockClientProtocol
func (_mock *MockClientProtocol[Req, Resp]) ParseResponse(data []byte) (Resp, error) {
	ret := _mock.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for ParseResponse")
	}

	var r0 Resp
	var r1 error
	if returnFunc, ok := ret.Get(0).(func([]byte) (Resp, error)); ok {
		return returnFunc(data)
	}
	if returnFunc, ok := ret.Get(0).(func([]byte) Resp); ok {
		r0 = returnFunc(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Resp)
		}
	}
	if returnFunc, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = returnFunc(data)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockClientProtocol_ParseResponse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseResponse'
type MockClientProtocol_ParseResponse_Call[Req tcpServer.Request, Resp tcpServer.Response] struct {
	*mock.Call
}

// ParseResponse is a helper method to define mock.On call
//   - data
func (_e *MockClientProtocol_Expecter[Req, Resp]) ParseResponse(data interface{}) *MockClientProtocol_ParseResponse_Call[Req, Resp] {
	return &MockClientProtocol_ParseResponse_Call[Req, Resp]{Call: _e.mock.On("ParseResponse", data)}
}

func (_c *MockClientProtocol_ParseResponse_Call[Req, Resp]) Run(run func(data []byte)) *MockClientProtocol_ParseResponse_Call[Req, Resp] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockClientProtocol_ParseResponse_Call[Req, Resp]) Return(v Resp, err error) *MockClientProtocol_ParseResponse_Call[Req, Resp] {
	_c.Call.Return(v, err)
	return _c
}

func (_c *MockClientProtocol_ParseResponse_Call[Req, Resp]) RunAndReturn(run func(data []byte) (Resp, error)) *MockClientProtocol_ParseResponse_Call[Req, Resp] {
	_c.Call.Return(run)
	return _c
}
