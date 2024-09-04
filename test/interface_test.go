package gobs_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xarest/gobs"
)

type TStruct struct {
	Name string
}

func (t *TStruct) Setup(ctx context.Context, deps ...gobs.IService) error {
	t.Name = deps[0].(string)
	return nil
}

type IStruct gobs.IServiceSetup

type InterfaceSuit struct {
	suite.Suite
}

func TestInterface(t *testing.T) {
	suite.Run(t, new(InterfaceSuit))
}

func (s *SchedulerSuit) Test_Struct_Struct() {
	deps := gobs.Dependencies{
		&TStruct{Name: "OK"},
	}

	t := &TStruct{}

	require.NoError(s.T(), deps.Assign(&t), "Assign expected no error")
	assert.Equal(s.T(), "OK", t.Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Struct_Interface() {
	deps := gobs.Dependencies{
		&TStruct{Name: "OK"},
	}

	var t IStruct

	require.NoError(s.T(), deps.Assign(&t), "Assign expected no error")
	assert.Equal(s.T(), "OK", t.(*TStruct).Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Interface_Struct() {
	var it IStruct = &TStruct{Name: "OK"}
	deps := gobs.Dependencies{it}

	ts := &TStruct{}
	require.NoError(s.T(), deps.Assign(&ts), "Assign expected no error")
	assert.Equal(s.T(), "OK", ts.Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Interface_Interface() {
	var it IStruct = &TStruct{Name: "OK"}
	deps := gobs.Dependencies{it}

	var t IStruct

	require.NoError(s.T(), deps.Assign(&t), "Assign expected no error")
	assert.Equal(s.T(), "OK", t.(*TStruct).Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Interface_Interface_Nil() {
	var it IStruct = &TStruct{Name: "OK"}
	deps := gobs.Dependencies{nil, it}

	var t IStruct

	require.NoError(s.T(), deps.Assign(nil, &t), "Assign expected no error")
	assert.Equal(s.T(), "OK", t.(*TStruct).Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Struct_Interface_Nil() {
	deps := gobs.Dependencies{nil, &TStruct{Name: "OK"}}

	var t IStruct

	require.NoError(s.T(), deps.Assign(nil, &t), "Assign expected no error")
	assert.Equal(s.T(), "OK", t.(*TStruct).Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Interface_Struct_Nil() {
	var it IStruct = &TStruct{Name: "OK"}
	deps := gobs.Dependencies{nil, it}

	ts := &TStruct{}
	require.NoError(s.T(), deps.Assign(nil, &ts), "Assign expected no error")
	assert.Equal(s.T(), "OK", ts.Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Struct_Struct_Nil() {
	deps := gobs.Dependencies{nil, &TStruct{Name: "OK"}}

	t := &TStruct{}

	require.NoError(s.T(), deps.Assign(nil, &t), "Assign expected no error")
	assert.Equal(s.T(), "OK", t.Name, "Expected t.Name to be OK")
}

func (s *SchedulerSuit) Test_Struct_Struct_And_Struct_Struct() {
	deps := gobs.Dependencies{
		&TStruct{Name: "OK"},
		&TStruct{Name: "OK2"},
	}

	var (
		t1 = &TStruct{}
		t2 = &TStruct{}
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Struct_Struct_And_Struct_Interface() {
	deps := gobs.Dependencies{
		&TStruct{Name: "OK"},
		&TStruct{Name: "OK2"},
	}

	var (
		t1 = &TStruct{}
		t2 IStruct
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.(*TStruct).Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Struct_Interface_And_Struct_Struct() {
	var it IStruct = &TStruct{Name: "OK2"}
	deps := gobs.Dependencies{
		&TStruct{Name: "OK"},
		it,
	}

	var (
		t1 = &TStruct{}
		t2 = &TStruct{}
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Struct_Interface_And_Struct_Interface() {
	var it IStruct = &TStruct{Name: "OK2"}
	deps := gobs.Dependencies{
		&TStruct{Name: "OK"},
		it,
	}

	var (
		t1 = &TStruct{}
		t2 IStruct
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.(*TStruct).Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Interface_Struct_And_Struct_Struct() {
	var it IStruct = &TStruct{Name: "OK"}
	deps := gobs.Dependencies{
		it,
		&TStruct{Name: "OK2"},
	}

	var (
		t1 = &TStruct{}
		t2 = &TStruct{}
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Interface_Struct_And_Struct_Interface() {
	var it IStruct = &TStruct{Name: "OK"}
	deps := gobs.Dependencies{
		it,
		&TStruct{Name: "OK2"},
	}

	var (
		t1 = &TStruct{}
		t2 IStruct
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.(*TStruct).Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Interface_Interface_And_Struct_Struct() {
	var it1 IStruct = &TStruct{Name: "OK"}
	var it2 IStruct = &TStruct{Name: "OK2"}
	deps := gobs.Dependencies{it1, it2}

	var (
		t1 = &TStruct{}
		t2 = &TStruct{}
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Interface_Interface_And_Struct_Interface() {
	var it1 IStruct = &TStruct{Name: "OK"}
	var it2 IStruct = &TStruct{Name: "OK2"}
	deps := gobs.Dependencies{
		it1, it2,
	}

	var (
		t1 = &TStruct{}
		t2 IStruct
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.(*TStruct).Name, "Expected t.Name to be OK2")
}

func (s *SchedulerSuit) Test_Interface_Interface_And_Interface_Interface() {
	var it1 IStruct = &TStruct{Name: "OK"}
	var it2 IStruct = &TStruct{Name: "OK2"}
	deps := gobs.Dependencies{
		it1,
		it2,
	}

	var (
		t1 = &TStruct{}
		t2 = &TStruct{}
	)

	require.NoError(s.T(), deps.Assign(&t1, &t2), "Assign expected no error")
	assert.Equal(s.T(), "OK", t1.Name, "Expected t.Name to be OK")
	assert.Equal(s.T(), "OK2", t2.Name, "Expected t.Name to be OK2")
}
