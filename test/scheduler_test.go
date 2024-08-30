package gobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
	"github.com/xarest/gobs/logger"
	"github.com/xarest/gobs/scheduler"
	"github.com/xarest/gobs/types"
)

type SchedulerSuit struct {
	suite.Suite
}

func TestScheduler(t *testing.T) {
	suite.Run(t, new(SchedulerSuit))
}

// func (s *SchedulerSuit) SetupSuite() {
// 	fmt.Println("SchedulerSuit/SetupSuite")
// }

// func (s *SchedulerSuit) TearDownSuite() {
// 	fmt.Println("SchedulerSuit/TearDownSuite")
// }

// func (s *SchedulerSuit) SetupTest() {
// 	fmt.Println("SchedulerSuit/SetupTest")
// }

// func (s *SchedulerSuit) TearDownTest() {
// 	fmt.Println("SchedulerSuit/TearDownTest")
// }

type MockTask struct {
	following []*MockTask
	followers []*MockTask
	name      string
	isAsync   bool
	delay     time.Duration
}

func (m *MockTask) Run(ctx context.Context, status common.ServiceStatus) error {
	time.Sleep(m.delay)
	return nil
}
func (m *MockTask) DependOn(status common.ServiceStatus) []types.ITask {
	out := make([]types.ITask, 0, len(m.following))
	for _, f := range m.following {
		out = append(out, f)
	}
	return out
}

func (m *MockTask) Followers(status common.ServiceStatus) []types.ITask {
	out := make([]types.ITask, 0, len(m.followers))
	for _, f := range m.followers {
		out = append(out, f)
	}
	return out
}

func (m *MockTask) IsRunAsync(status common.ServiceStatus) bool {
	return m.isAsync
}

func (m *MockTask) Name() string {
	return m.name
}

var _ types.ITask = &MockTask{}

func (s *SchedulerSuit) TestSchedulerSync() {
	var (
		taskA = MockTask{name: "D"}
		taskB = MockTask{name: "C", following: []*MockTask{&taskA}}
		taskC = MockTask{name: "B", following: []*MockTask{&taskB}}
		taskD = MockTask{name: "A", following: []*MockTask{&taskB, &taskA}}
	)
	taskA.followers = []*MockTask{&taskB, &taskD}
	taskB.followers = []*MockTask{&taskC, &taskD}
	taskC.followers = []*MockTask{&taskD}
	scheduler := scheduler.NewScheduler(
		context.TODO(),
		logger.NewLog(nil),
		[]types.ITask{
			&taskA, &taskB, &taskC, &taskD,
		},
		common.StatusSetup,
		gobs.DEFAULT_MAX_CONCURRENT,
	)
	err := scheduler.Run(context.TODO())
	require.Nil(s.T(), err)
	results, err := scheduler.Release()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 4, len(results))
	assert.Equal(s.T(), "D", results[0].Name())
	assert.Equal(s.T(), "C", results[1].Name())
	assert.Equal(s.T(), "B", results[2].Name())
	assert.Equal(s.T(), "A", results[3].Name())
}

func (s *SchedulerSuit) TestSchedulerAsyncWithoutDelay() {
	var (
		taskA = MockTask{name: "A"}
		taskB = MockTask{name: "B", following: []*MockTask{&taskA}}
		taskC = MockTask{name: "C", following: []*MockTask{&taskB}}
		taskD = MockTask{name: "D", following: []*MockTask{&taskB, &taskA}}
	)
	taskA.followers = []*MockTask{&taskB, &taskD}
	taskB.followers = []*MockTask{&taskC, &taskD}
	taskC.followers = []*MockTask{&taskD}
	scheduler := scheduler.NewScheduler(
		context.TODO(),
		logger.NewLog(nil),
		[]types.ITask{
			&taskA, &taskB, &taskC, &taskD,
		},
		common.StatusSetup,
		gobs.DEFAULT_MAX_CONCURRENT,
	)
	err := scheduler.Run(context.TODO())
	require.Nil(s.T(), err)
	results, err := scheduler.Release()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 4, len(results))
	assert.Equal(s.T(), "A", results[0].Name())
	assert.Equal(s.T(), "B", results[1].Name())
	assert.Equal(s.T(), "C", results[2].Name())
	assert.Equal(s.T(), "D", results[3].Name())
}

func (s *SchedulerSuit) TestSchedulerAsyncWithDelay() {
	var (
		taskE = MockTask{name: "E"}
		taskA = MockTask{name: "A", following: []*MockTask{&taskE}}
		taskB = MockTask{name: "B", following: []*MockTask{&taskA, &taskE}}
		taskC = MockTask{name: "C", following: []*MockTask{&taskB}, isAsync: true, delay: 1 * time.Second}
		taskD = MockTask{name: "D", following: []*MockTask{&taskB, &taskA}}
	)
	taskE.followers = []*MockTask{&taskA, &taskB}
	taskA.followers = []*MockTask{&taskB, &taskD}
	taskB.followers = []*MockTask{&taskC, &taskD}
	taskC.followers = []*MockTask{&taskD}
	scheduler := scheduler.NewScheduler(
		context.TODO(),
		logger.NewLog(nil),
		[]types.ITask{
			&taskA, &taskB, &taskC, &taskD, &taskE,
		},
		common.StatusSetup,
		gobs.DEFAULT_MAX_CONCURRENT,
	)
	err := scheduler.Run(context.TODO())
	require.Nil(s.T(), err)
	scheduler.Interrupt()
	results, err := scheduler.Release()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 5, len(results))
	assert.Equal(s.T(), "E", results[0].Name())
	assert.Equal(s.T(), "A", results[1].Name())
	assert.Equal(s.T(), "B", results[2].Name())
	assert.Equal(s.T(), "D", results[3].Name())
	assert.Equal(s.T(), "C", results[4].Name())
}

func (s *SchedulerSuit) TestSchedulerWithAsyncInterruptAndStop() {
	var (
		taskE = MockTask{name: "E"}
		taskA = MockTask{name: "A", following: []*MockTask{&taskE}}
		taskB = MockTask{name: "B", following: []*MockTask{&taskA, &taskE}}
		taskC = MockTask{name: "C", following: []*MockTask{&taskB}, isAsync: true, delay: 1 * time.Second}
		taskD = MockTask{name: "D", following: []*MockTask{&taskB, &taskA}}
	)
	taskE.followers = []*MockTask{&taskA, &taskB}
	taskA.followers = []*MockTask{&taskB, &taskD}
	taskB.followers = []*MockTask{&taskC, &taskD}
	taskC.followers = []*MockTask{&taskD}
	sched := scheduler.NewScheduler(
		context.TODO(),
		logger.NewLog(nil),
		[]types.ITask{
			&taskA, &taskB, &taskC, &taskD, &taskE,
		},
		common.StatusSetup,
		gobs.DEFAULT_MAX_CONCURRENT,
	)

	go sched.Run(context.TODO())
	time.Sleep(500 * time.Millisecond)
	sched.Interrupt()
	results, err := sched.Release()
	assert.Equal(s.T(), context.Canceled, err)
	require.Equal(s.T(), 4, len(results))
	assert.Equal(s.T(), "E", results[0].Name())
	assert.Equal(s.T(), "A", results[1].Name())
	assert.Equal(s.T(), "B", results[2].Name())
	assert.Equal(s.T(), "D", results[3].Name())

	sched = scheduler.NewScheduler(
		context.TODO(),
		logger.NewLog(nil),
		results,
		common.StatusStop,
		gobs.DEFAULT_MAX_CONCURRENT,
	)

	go sched.Run(context.TODO())
	time.Sleep(500 * time.Millisecond)
	results, err = sched.Release()
	assert.Nil(s.T(), err)
	require.Equal(s.T(), 4, len(results))
	assert.Equal(s.T(), "E", results[0].Name())
	assert.Equal(s.T(), "A", results[1].Name())
	assert.Equal(s.T(), "B", results[2].Name())
	assert.Equal(s.T(), "D", results[3].Name())
}
