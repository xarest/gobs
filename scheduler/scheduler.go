package scheduler

import (
	"context"
	"sync"

	"github.com/traphamxuan/gobs/common"
	"github.com/traphamxuan/gobs/logger"
	"github.com/traphamxuan/gobs/types"
	"github.com/traphamxuan/gobs/utils"
)

type Scheduler struct {
	*logger.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	status        common.ServiceStatus
	wg            sync.WaitGroup
	chReqSync     chan types.ITask
	chReqAsync    chan types.ITask
	chRes         chan types.ITask
	chErr         chan error
	isRunning     map[string]bool
	isFinished    map[string]bool
	mutexRun      *sync.RWMutex
	mutexFinished *sync.RWMutex
	err           error
	ranList       []types.ITask
	finishedList  []types.ITask
	Tasks         []types.ITask
}

func NewScheduler(
	ctx context.Context,
	log *logger.Logger,
	tasks []types.ITask,
	ss common.ServiceStatus) *Scheduler {
	ctx, cancel := context.WithCancel(ctx)
	log.AddTag("Scheduler/" + ss.String())
	return &Scheduler{
		Logger:        log,
		ctx:           ctx,
		cancel:        cancel,
		status:        ss,
		chReqSync:     make(chan types.ITask, len(tasks)),
		chReqAsync:    make(chan types.ITask, len(tasks)),
		chRes:         make(chan types.ITask, len(tasks)),
		chErr:         make(chan error, len(tasks)),
		ranList:       make([]types.ITask, 0, len(tasks)),
		finishedList:  make([]types.ITask, 0, len(tasks)),
		isRunning:     make(map[string]bool, len(tasks)),
		isFinished:    make(map[string]bool, len(tasks)),
		mutexRun:      &sync.RWMutex{},
		mutexFinished: &sync.RWMutex{},
		Tasks:         tasks,
		err:           nil,
	}
}

func (r *Scheduler) Interrupt() {
	r.cancel()
}

func (r *Scheduler) Release() ([]types.ITask, error) {
	r.wg.Wait()
	return r.finishedList, r.err
}

func (r *Scheduler) RunSync(ctx context.Context) ([]types.ITask, error) {
	r.wg.Add(1)
	defer r.wg.Done()
	r.err = r.startSyncRun(ctx, r.Tasks)
	return r.finishedList, r.err
}

func (r *Scheduler) startSyncRun(ctx context.Context, tasks []types.ITask) error {
	for _, task := range tasks {
		if r.ctx.Err() != nil {
			return r.ctx.Err()
		}
		key := task.Name()

		r.mutexFinished.RLock()
		if isFinished, ok := r.isFinished[key]; !ok || !isFinished {
			r.mutexFinished.RUnlock()
			if err := r.startSyncRun(ctx, task.DependOn(r.status)); err != nil {
				return err
			}
			if err := task.Run(ctx, r.status); utils.WrapCommonError(err) != nil {
				return err
			}
			r.mutexFinished.Lock()
			r.isFinished[key] = true
			r.mutexFinished.Unlock()
			r.finishedList = append(r.finishedList, task)
		} else {
			r.mutexFinished.RUnlock()
		}
	}
	return nil
}

func (r *Scheduler) RunAsync(ctx context.Context) ([]types.ITask, error) {
	r.wg.Add(1)
	defer r.wg.Done()
	go r.startProducer()
	go r.startConsumer(ctx)

	r.checkAndLoad(r.Tasks, func(task types.ITask) {
		if task.IsRunAsync(r.status) {
			r.chReqAsync <- task
		} else {
			r.chReqAsync <- task
		}
	})
	select {
	case <-r.ctx.Done():
		r.err = r.ctx.Err()
	case err := <-r.chErr:
		r.err = err
	}
	return r.finishedList, r.err
}

func (r *Scheduler) startProducer() {
	defer func() {
		close(r.chReqSync)
		close(r.chReqAsync)
	}()
	utils.WaitOnEvents(r.ctx, func(ctx context.Context, task types.ITask) error {
		key := task.Name()
		r.finishedList = append(r.finishedList, task)
		r.mutexFinished.Lock()
		r.isFinished[key] = true
		r.mutexFinished.Unlock()
		if len(r.finishedList) == len(r.Tasks) {
			return utils.ErrorEndOfProcessing
		}
		followers := task.Followers(r.status)
		r.checkAndLoad(followers, func(task types.ITask) {
			if task.IsRunAsync(r.status) {
				r.chReqAsync <- task
			} else {
				r.chReqSync <- task
			}
		})
		return nil
	}, r.chErr, r.chRes)
}

func (r *Scheduler) startConsumer(ctx context.Context) {
	defer func() {
		close(r.chRes)
		close(r.chErr)
	}()
	var wgTask sync.WaitGroup
	wgTask.Add(1)
	go r.startSyncWorker(ctx, &wgTask)
	wgTask.Add(1)
	go r.startAsyncWorker(ctx, &wgTask)
	wgTask.Wait()
}

func (r *Scheduler) startSyncWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.WaitOnEvents(r.ctx, func(_ context.Context, task types.ITask) error {
		if err := task.Run(ctx, r.status); utils.WrapCommonError(err) != nil {
			r.chErr <- err
			return err
		}
		r.chRes <- task
		return nil
	}, r.chErr, r.chReqSync)
}

func (r *Scheduler) startAsyncWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.WaitOnEvents(r.ctx, func(_ context.Context, task types.ITask) error {
		wg.Add(1)
		go func(task types.ITask) {
			defer wg.Done()
			if err := task.Run(ctx, r.status); utils.WrapCommonError(err) != nil {
				r.chErr <- err
				return
			}
			r.chRes <- task
		}(task)
		return nil
	}, r.chErr, r.chReqAsync)
}

func (r *Scheduler) checkAndLoad(tasks []types.ITask, onLoad func(task types.ITask)) {
	for _, task := range tasks {
		// Check if task's dependencies finished
		if r.checkDependenciesReady(task) {
			key := task.Name()
			// Only push to channel if the task is not running
			r.mutexRun.RLock()
			if isRunning, ok := r.isRunning[key]; !ok || !isRunning {
				r.mutexRun.RUnlock()
				r.mutexRun.Lock()
				r.isRunning[key] = true
				r.mutexRun.Unlock()
				r.ranList = append(r.ranList, task)
				onLoad(task)
			} else {
				r.mutexRun.RUnlock()
			}
		}
	}
}

func (r *Scheduler) checkDependenciesReady(task types.ITask) bool {
	r.mutexFinished.RLock()
	defer r.mutexFinished.RUnlock()
	for _, dep := range task.DependOn(r.status) {
		if isFinished, ok := r.isFinished[dep.Name()]; !ok || !isFinished {
			return false
		}
	}
	return true
}
