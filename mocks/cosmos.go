// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/umee-network/peggo/cmd/peggo/client (interfaces: CosmosClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	client "github.com/cosmos/cosmos-sdk/client"
	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockCosmosClient is a mock of CosmosClient interface.
type MockCosmosClient struct {
	ctrl     *gomock.Controller
	recorder *MockCosmosClientMockRecorder
}

// MockCosmosClientMockRecorder is the mock recorder for MockCosmosClient.
type MockCosmosClientMockRecorder struct {
	mock *MockCosmosClient
}

// NewMockCosmosClient creates a new mock instance.
func NewMockCosmosClient(ctrl *gomock.Controller) *MockCosmosClient {
	mock := &MockCosmosClient{ctrl: ctrl}
	mock.recorder = &MockCosmosClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCosmosClient) EXPECT() *MockCosmosClientMockRecorder {
	return m.recorder
}

// AsyncBroadcastMsg mocks base method.
func (m *MockCosmosClient) AsyncBroadcastMsg(arg0 ...types.Msg) (*types.TxResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AsyncBroadcastMsg", varargs...)
	ret0, _ := ret[0].(*types.TxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AsyncBroadcastMsg indicates an expected call of AsyncBroadcastMsg.
func (mr *MockCosmosClientMockRecorder) AsyncBroadcastMsg(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsyncBroadcastMsg", reflect.TypeOf((*MockCosmosClient)(nil).AsyncBroadcastMsg), arg0...)
}

// CanSignTransactions mocks base method.
func (m *MockCosmosClient) CanSignTransactions() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanSignTransactions")
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanSignTransactions indicates an expected call of CanSignTransactions.
func (mr *MockCosmosClientMockRecorder) CanSignTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanSignTransactions", reflect.TypeOf((*MockCosmosClient)(nil).CanSignTransactions))
}

// ClientContext mocks base method.
func (m *MockCosmosClient) ClientContext() client.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientContext")
	ret0, _ := ret[0].(client.Context)
	return ret0
}

// ClientContext indicates an expected call of ClientContext.
func (mr *MockCosmosClientMockRecorder) ClientContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientContext", reflect.TypeOf((*MockCosmosClient)(nil).ClientContext))
}

// Close mocks base method.
func (m *MockCosmosClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockCosmosClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCosmosClient)(nil).Close))
}

// FromAddress mocks base method.
func (m *MockCosmosClient) FromAddress() types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromAddress")
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// FromAddress indicates an expected call of FromAddress.
func (mr *MockCosmosClientMockRecorder) FromAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromAddress", reflect.TypeOf((*MockCosmosClient)(nil).FromAddress))
}

// QueryClient mocks base method.
func (m *MockCosmosClient) QueryClient() *grpc.ClientConn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryClient")
	ret0, _ := ret[0].(*grpc.ClientConn)
	return ret0
}

// QueryClient indicates an expected call of QueryClient.
func (mr *MockCosmosClientMockRecorder) QueryClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryClient", reflect.TypeOf((*MockCosmosClient)(nil).QueryClient))
}

// QueueBroadcastMsg mocks base method.
func (m *MockCosmosClient) QueueBroadcastMsg(arg0 ...types.Msg) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueueBroadcastMsg", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueueBroadcastMsg indicates an expected call of QueueBroadcastMsg.
func (mr *MockCosmosClientMockRecorder) QueueBroadcastMsg(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueBroadcastMsg", reflect.TypeOf((*MockCosmosClient)(nil).QueueBroadcastMsg), arg0...)
}

// SyncBroadcastMsg mocks base method.
func (m *MockCosmosClient) SyncBroadcastMsg(arg0 ...types.Msg) (*types.TxResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SyncBroadcastMsg", varargs...)
	ret0, _ := ret[0].(*types.TxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncBroadcastMsg indicates an expected call of SyncBroadcastMsg.
func (mr *MockCosmosClientMockRecorder) SyncBroadcastMsg(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncBroadcastMsg", reflect.TypeOf((*MockCosmosClient)(nil).SyncBroadcastMsg), arg0...)
}
