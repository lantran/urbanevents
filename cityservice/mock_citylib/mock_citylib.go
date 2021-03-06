// Automatically generated by MockGen. DO NOT EDIT!
// Source: citylib/instagramlinkenricher.go

package mock_citylib

import (
	instagram "github.com/dimroc/go-instagram/instagram"
	gomock "github.com/golang/mock/gomock"
)

// Mock of MediaRetriever interface
type MockMediaRetriever struct {
	ctrl     *gomock.Controller
	recorder *_MockMediaRetrieverRecorder
}

// Recorder for MockMediaRetriever (not exported)
type _MockMediaRetrieverRecorder struct {
	mock *MockMediaRetriever
}

func NewMockMediaRetriever(ctrl *gomock.Controller) *MockMediaRetriever {
	mock := &MockMediaRetriever{ctrl: ctrl}
	mock.recorder = &_MockMediaRetrieverRecorder{mock}
	return mock
}

func (_m *MockMediaRetriever) EXPECT() *_MockMediaRetrieverRecorder {
	return _m.recorder
}

func (_m *MockMediaRetriever) GetShortcode(shortcode string) (*instagram.Media, error) {
	ret := _m.ctrl.Call(_m, "GetShortcode", shortcode)
	ret0, _ := ret[0].(*instagram.Media)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockMediaRetrieverRecorder) GetShortcode(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetShortcode", arg0)
}
