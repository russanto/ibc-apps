// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types (interfaces: ChannelKeeper)
//
// Generated by this command:
//
//	mockgen -package=mock -destination=./test/mock/channel_keeper.go github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types ChannelKeeper
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/ibc-go/modules/capability/types"
	types1 "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	gomock "go.uber.org/mock/gomock"
)

// MockChannelKeeper is a mock of ChannelKeeper interface.
type MockChannelKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockChannelKeeperMockRecorder
}

// MockChannelKeeperMockRecorder is the mock recorder for MockChannelKeeper.
type MockChannelKeeperMockRecorder struct {
	mock *MockChannelKeeper
}

// NewMockChannelKeeper creates a new mock instance.
func NewMockChannelKeeper(ctrl *gomock.Controller) *MockChannelKeeper {
	mock := &MockChannelKeeper{ctrl: ctrl}
	mock.recorder = &MockChannelKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChannelKeeper) EXPECT() *MockChannelKeeperMockRecorder {
	return m.recorder
}

// GetAllChannelsWithPortPrefix mocks base method.
func (m *MockChannelKeeper) GetAllChannelsWithPortPrefix(arg0 types.Context, arg1 string) []types1.IdentifiedChannel {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChannelsWithPortPrefix", arg0, arg1)
	ret0, _ := ret[0].([]types1.IdentifiedChannel)
	return ret0
}

// GetAllChannelsWithPortPrefix indicates an expected call of GetAllChannelsWithPortPrefix.
func (mr *MockChannelKeeperMockRecorder) GetAllChannelsWithPortPrefix(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChannelsWithPortPrefix", reflect.TypeOf((*MockChannelKeeper)(nil).GetAllChannelsWithPortPrefix), arg0, arg1)
}

// GetChannel mocks base method.
func (m *MockChannelKeeper) GetChannel(arg0 types.Context, arg1, arg2 string) (types1.Channel, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannel", arg0, arg1, arg2)
	ret0, _ := ret[0].(types1.Channel)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetChannel indicates an expected call of GetChannel.
func (mr *MockChannelKeeperMockRecorder) GetChannel(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannel", reflect.TypeOf((*MockChannelKeeper)(nil).GetChannel), arg0, arg1, arg2)
}

// GetNextSequenceSend mocks base method.
func (m *MockChannelKeeper) GetNextSequenceSend(arg0 types.Context, arg1, arg2 string) (uint64, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextSequenceSend", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetNextSequenceSend indicates an expected call of GetNextSequenceSend.
func (mr *MockChannelKeeperMockRecorder) GetNextSequenceSend(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextSequenceSend", reflect.TypeOf((*MockChannelKeeper)(nil).GetNextSequenceSend), arg0, arg1, arg2)
}

// GetPacketCommitment mocks base method.
func (m *MockChannelKeeper) GetPacketCommitment(arg0 types.Context, arg1, arg2 string, arg3 uint64) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPacketCommitment", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetPacketCommitment indicates an expected call of GetPacketCommitment.
func (mr *MockChannelKeeperMockRecorder) GetPacketCommitment(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPacketCommitment", reflect.TypeOf((*MockChannelKeeper)(nil).GetPacketCommitment), arg0, arg1, arg2, arg3)
}

// LookupModuleByChannel mocks base method.
func (m *MockChannelKeeper) LookupModuleByChannel(arg0 types.Context, arg1, arg2 string) (string, *types0.Capability, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupModuleByChannel", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*types0.Capability)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LookupModuleByChannel indicates an expected call of LookupModuleByChannel.
func (mr *MockChannelKeeperMockRecorder) LookupModuleByChannel(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupModuleByChannel", reflect.TypeOf((*MockChannelKeeper)(nil).LookupModuleByChannel), arg0, arg1, arg2)
}
