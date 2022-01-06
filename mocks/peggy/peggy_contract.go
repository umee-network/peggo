// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/umee-network/peggo/orchestrator/ethereum/peggy (interfaces: Contract)

// Package peggy is a generated GoMock package.
package peggy

import (
	context "context"
	big "math/big"
	reflect "reflect"
	time "time"

	types "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	common "github.com/ethereum/go-ethereum/common"
	gomock "github.com/golang/mock/gomock"
	peggy "github.com/umee-network/peggo/orchestrator/ethereum/peggy"
	provider "github.com/umee-network/peggo/orchestrator/ethereum/provider"
)

// MockContract is a mock of Contract interface.
type MockContract struct {
	ctrl     *gomock.Controller
	recorder *MockContractMockRecorder
}

// MockContractMockRecorder is the mock recorder for MockContract.
type MockContractMockRecorder struct {
	mock *MockContract
}

// NewMockContract creates a new mock instance.
func NewMockContract(ctrl *gomock.Controller) *MockContract {
	mock := &MockContract{ctrl: ctrl}
	mock.recorder = &MockContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContract) EXPECT() *MockContractMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockContract) Address() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockContractMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockContract)(nil).Address))
}

// EncodeTransactionBatch mocks base method.
func (m *MockContract) EncodeTransactionBatch(arg0 context.Context, arg1 *types.Valset, arg2 *types.OutgoingTxBatch, arg3 []*types.MsgConfirmBatch) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncodeTransactionBatch", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncodeTransactionBatch indicates an expected call of EncodeTransactionBatch.
func (mr *MockContractMockRecorder) EncodeTransactionBatch(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodeTransactionBatch", reflect.TypeOf((*MockContract)(nil).EncodeTransactionBatch), arg0, arg1, arg2, arg3)
}

// EncodeValsetUpdate mocks base method.
func (m *MockContract) EncodeValsetUpdate(arg0 context.Context, arg1, arg2 *types.Valset, arg3 []*types.MsgValsetConfirm) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncodeValsetUpdate", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncodeValsetUpdate indicates an expected call of EncodeValsetUpdate.
func (mr *MockContractMockRecorder) EncodeValsetUpdate(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodeValsetUpdate", reflect.TypeOf((*MockContract)(nil).EncodeValsetUpdate), arg0, arg1, arg2, arg3)
}

// EstimateGas mocks base method.
func (m *MockContract) EstimateGas(arg0 context.Context, arg1 common.Address, arg2 []byte) (uint64, *big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateGas", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(*big.Int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EstimateGas indicates an expected call of EstimateGas.
func (mr *MockContractMockRecorder) EstimateGas(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateGas", reflect.TypeOf((*MockContract)(nil).EstimateGas), arg0, arg1, arg2)
}

// FromAddress mocks base method.
func (m *MockContract) FromAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// FromAddress indicates an expected call of FromAddress.
func (mr *MockContractMockRecorder) FromAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromAddress", reflect.TypeOf((*MockContract)(nil).FromAddress))
}

// GetERC20Decimals mocks base method.
func (m *MockContract) GetERC20Decimals(arg0 context.Context, arg1, arg2 common.Address) (byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetERC20Decimals", arg0, arg1, arg2)
	ret0, _ := ret[0].(byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetERC20Decimals indicates an expected call of GetERC20Decimals.
func (mr *MockContractMockRecorder) GetERC20Decimals(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetERC20Decimals", reflect.TypeOf((*MockContract)(nil).GetERC20Decimals), arg0, arg1, arg2)
}

// GetERC20Symbol mocks base method.
func (m *MockContract) GetERC20Symbol(arg0 context.Context, arg1, arg2 common.Address) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetERC20Symbol", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetERC20Symbol indicates an expected call of GetERC20Symbol.
func (mr *MockContractMockRecorder) GetERC20Symbol(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetERC20Symbol", reflect.TypeOf((*MockContract)(nil).GetERC20Symbol), arg0, arg1, arg2)
}

// GetPeggyID mocks base method.
func (m *MockContract) GetPeggyID(arg0 context.Context, arg1 common.Address) (common.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeggyID", arg0, arg1)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeggyID indicates an expected call of GetPeggyID.
func (mr *MockContractMockRecorder) GetPeggyID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeggyID", reflect.TypeOf((*MockContract)(nil).GetPeggyID), arg0, arg1)
}

// GetPendingTxInputList mocks base method.
func (m *MockContract) GetPendingTxInputList() *peggy.PendingTxInputList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingTxInputList")
	ret0, _ := ret[0].(*peggy.PendingTxInputList)
	return ret0
}

// GetPendingTxInputList indicates an expected call of GetPendingTxInputList.
func (mr *MockContractMockRecorder) GetPendingTxInputList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingTxInputList", reflect.TypeOf((*MockContract)(nil).GetPendingTxInputList))
}

// GetTxBatchNonce mocks base method.
func (m *MockContract) GetTxBatchNonce(arg0 context.Context, arg1, arg2 common.Address) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxBatchNonce", arg0, arg1, arg2)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxBatchNonce indicates an expected call of GetTxBatchNonce.
func (mr *MockContractMockRecorder) GetTxBatchNonce(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxBatchNonce", reflect.TypeOf((*MockContract)(nil).GetTxBatchNonce), arg0, arg1, arg2)
}

// GetValsetNonce mocks base method.
func (m *MockContract) GetValsetNonce(arg0 context.Context, arg1 common.Address) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValsetNonce", arg0, arg1)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValsetNonce indicates an expected call of GetValsetNonce.
func (mr *MockContractMockRecorder) GetValsetNonce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValsetNonce", reflect.TypeOf((*MockContract)(nil).GetValsetNonce), arg0, arg1)
}

// IsPendingTxInput mocks base method.
func (m *MockContract) IsPendingTxInput(arg0 []byte, arg1 time.Duration) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPendingTxInput", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPendingTxInput indicates an expected call of IsPendingTxInput.
func (mr *MockContractMockRecorder) IsPendingTxInput(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPendingTxInput", reflect.TypeOf((*MockContract)(nil).IsPendingTxInput), arg0, arg1)
}

// Provider mocks base method.
func (m *MockContract) Provider() provider.EVMProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Provider")
	ret0, _ := ret[0].(provider.EVMProvider)
	return ret0
}

// Provider indicates an expected call of Provider.
func (mr *MockContractMockRecorder) Provider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Provider", reflect.TypeOf((*MockContract)(nil).Provider))
}

// SendTx mocks base method.
func (m *MockContract) SendTx(arg0 context.Context, arg1 common.Address, arg2 []byte, arg3 uint64, arg4 *big.Int) (common.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTx", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTx indicates an expected call of SendTx.
func (mr *MockContractMockRecorder) SendTx(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTx", reflect.TypeOf((*MockContract)(nil).SendTx), arg0, arg1, arg2, arg3, arg4)
}

// SubscribeToPendingTxs mocks base method.
func (m *MockContract) SubscribeToPendingTxs(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToPendingTxs", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeToPendingTxs indicates an expected call of SubscribeToPendingTxs.
func (mr *MockContractMockRecorder) SubscribeToPendingTxs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToPendingTxs", reflect.TypeOf((*MockContract)(nil).SubscribeToPendingTxs), arg0, arg1)
}
