// Code generated by MockGen. DO NOT EDIT.
// Source: supply.go

// Package supply_test is a generated GoMock package.
package supply_test

import (
	io "io"
	exec "os/exec"
	reflect "reflect"

	libbuildpack "github.com/cloudfoundry/libbuildpack"
	gomock "github.com/golang/mock/gomock"
)

// MockCommand is a mock of Command interface.
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand.
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance.
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockCommand) Execute(arg0 string, arg1, arg2 io.Writer, arg3 string, arg4 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockCommandMockRecorder) Execute(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCommand)(nil).Execute), varargs...)
}

// Output mocks base method.
func (m *MockCommand) Output(arg0, arg1 string, arg2 ...string) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Output", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Output indicates an expected call of Output.
func (mr *MockCommandMockRecorder) Output(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockCommand)(nil).Output), varargs...)
}

// Run mocks base method.
func (m *MockCommand) Run(cmd *exec.Cmd) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", cmd)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockCommandMockRecorder) Run(cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCommand)(nil).Run), cmd)
}

// RunWithOutput mocks base method.
func (m *MockCommand) RunWithOutput(cmd *exec.Cmd) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunWithOutput", cmd)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunWithOutput indicates an expected call of RunWithOutput.
func (mr *MockCommandMockRecorder) RunWithOutput(cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunWithOutput", reflect.TypeOf((*MockCommand)(nil).RunWithOutput), cmd)
}

// MockManifest is a mock of Manifest interface.
type MockManifest struct {
	ctrl     *gomock.Controller
	recorder *MockManifestMockRecorder
}

// MockManifestMockRecorder is the mock recorder for MockManifest.
type MockManifestMockRecorder struct {
	mock *MockManifest
}

// NewMockManifest creates a new mock instance.
func NewMockManifest(ctrl *gomock.Controller) *MockManifest {
	mock := &MockManifest{ctrl: ctrl}
	mock.recorder = &MockManifestMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManifest) EXPECT() *MockManifestMockRecorder {
	return m.recorder
}

// AllDependencyVersions mocks base method.
func (m *MockManifest) AllDependencyVersions(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllDependencyVersions", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// AllDependencyVersions indicates an expected call of AllDependencyVersions.
func (mr *MockManifestMockRecorder) AllDependencyVersions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllDependencyVersions", reflect.TypeOf((*MockManifest)(nil).AllDependencyVersions), arg0)
}

// DefaultVersion mocks base method.
func (m *MockManifest) DefaultVersion(depName string) (libbuildpack.Dependency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultVersion", depName)
	ret0, _ := ret[0].(libbuildpack.Dependency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DefaultVersion indicates an expected call of DefaultVersion.
func (mr *MockManifestMockRecorder) DefaultVersion(depName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultVersion", reflect.TypeOf((*MockManifest)(nil).DefaultVersion), depName)
}

// RootDir mocks base method.
func (m *MockManifest) RootDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RootDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// RootDir indicates an expected call of RootDir.
func (mr *MockManifestMockRecorder) RootDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RootDir", reflect.TypeOf((*MockManifest)(nil).RootDir))
}

// MockInstaller is a mock of Installer interface.
type MockInstaller struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerMockRecorder
}

// MockInstallerMockRecorder is the mock recorder for MockInstaller.
type MockInstallerMockRecorder struct {
	mock *MockInstaller
}

// NewMockInstaller creates a new mock instance.
func NewMockInstaller(ctrl *gomock.Controller) *MockInstaller {
	mock := &MockInstaller{ctrl: ctrl}
	mock.recorder = &MockInstallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstaller) EXPECT() *MockInstallerMockRecorder {
	return m.recorder
}

// InstallDependency mocks base method.
func (m *MockInstaller) InstallDependency(dep libbuildpack.Dependency, outputDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallDependency", dep, outputDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// InstallDependency indicates an expected call of InstallDependency.
func (mr *MockInstallerMockRecorder) InstallDependency(dep, outputDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallDependency", reflect.TypeOf((*MockInstaller)(nil).InstallDependency), dep, outputDir)
}

// InstallOnlyVersion mocks base method.
func (m *MockInstaller) InstallOnlyVersion(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallOnlyVersion", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InstallOnlyVersion indicates an expected call of InstallOnlyVersion.
func (mr *MockInstallerMockRecorder) InstallOnlyVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallOnlyVersion", reflect.TypeOf((*MockInstaller)(nil).InstallOnlyVersion), arg0, arg1)
}

// MockStager is a mock of Stager interface.
type MockStager struct {
	ctrl     *gomock.Controller
	recorder *MockStagerMockRecorder
}

// MockStagerMockRecorder is the mock recorder for MockStager.
type MockStagerMockRecorder struct {
	mock *MockStager
}

// NewMockStager creates a new mock instance.
func NewMockStager(ctrl *gomock.Controller) *MockStager {
	mock := &MockStager{ctrl: ctrl}
	mock.recorder = &MockStagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStager) EXPECT() *MockStagerMockRecorder {
	return m.recorder
}

// AddBinDependencyLink mocks base method.
func (m *MockStager) AddBinDependencyLink(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBinDependencyLink", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBinDependencyLink indicates an expected call of AddBinDependencyLink.
func (mr *MockStagerMockRecorder) AddBinDependencyLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBinDependencyLink", reflect.TypeOf((*MockStager)(nil).AddBinDependencyLink), arg0, arg1)
}

// BuildDir mocks base method.
func (m *MockStager) BuildDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// BuildDir indicates an expected call of BuildDir.
func (mr *MockStagerMockRecorder) BuildDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildDir", reflect.TypeOf((*MockStager)(nil).BuildDir))
}

// DepDir mocks base method.
func (m *MockStager) DepDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// DepDir indicates an expected call of DepDir.
func (mr *MockStagerMockRecorder) DepDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepDir", reflect.TypeOf((*MockStager)(nil).DepDir))
}

// DepsDir mocks base method.
func (m *MockStager) DepsDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepsDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// DepsDir indicates an expected call of DepsDir.
func (mr *MockStagerMockRecorder) DepsDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepsDir", reflect.TypeOf((*MockStager)(nil).DepsDir))
}

// DepsIdx mocks base method.
func (m *MockStager) DepsIdx() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepsIdx")
	ret0, _ := ret[0].(string)
	return ret0
}

// DepsIdx indicates an expected call of DepsIdx.
func (mr *MockStagerMockRecorder) DepsIdx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepsIdx", reflect.TypeOf((*MockStager)(nil).DepsIdx))
}

// WriteProfileD mocks base method.
func (m *MockStager) WriteProfileD(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteProfileD", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteProfileD indicates an expected call of WriteProfileD.
func (mr *MockStagerMockRecorder) WriteProfileD(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteProfileD", reflect.TypeOf((*MockStager)(nil).WriteProfileD), arg0, arg1)
}
