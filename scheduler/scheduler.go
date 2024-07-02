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
	ctx                context.Context
	cancel             context.CancelFunc
	status             common.ServiceStatus
	wg                 sync.WaitGroup
	chReqSync          chan types.ITask
	chReqAsync         chan types.ITask
	chRes              chan types.ITask
	chErr              chan error
	numOfConcurrencies int
	isRunning          map[string]bool
	isFinished         map[string]bool
	mutexRun           *sync.RWMutex
	mutexFinished      *sync.RWMutex
	err                error
	ranList            []types.ITask
	finishedList       []types.ITask
	Tasks              []types.ITask
}

func NewScheduler(
	ctx context.Context,
	log *logger.Logger,
	tasks []types.ITask,
	ss common.ServiceStatus,
	numOfConcurrencies int,
) *Scheduler {
	numOfTasks := len(tasks)
	concurrentLimit := numOfConcurrencies
	if concurrentLimit < 0 || concurrentLimit > numOfTasks {
		concurrentLimit = numOfTasks
	}
	ctx, cancel := context.WithCancel(ctx)
	log.AddTag("Scheduler-" + ss.String())
	sched := &Scheduler{
		Logger:             log,
		ctx:                ctx,
		cancel:             cancel,
		status:             ss,
		chReqSync:          make(chan types.ITask, numOfTasks),
		chReqAsync:         make(chan types.ITask, concurrentLimit),
		chRes:              make(chan types.ITask, numOfTasks),
		chErr:              make(chan error, numOfTasks),
		numOfConcurrencies: numOfConcurrencies,
		ranList:            make([]types.ITask, 0, numOfTasks),
		finishedList:       make([]types.ITask, 0, numOfTasks),
		isRunning:          make(map[string]bool, numOfTasks),
		isFinished:         make(map[string]bool, numOfTasks),
		mutexRun:           &sync.RWMutex{},
		mutexFinished:      &sync.RWMutex{},
		Tasks:              tasks,
		err:                nil,
	}
	return sched
}

func (r *Scheduler) SetIgnore(t types.ITask) {
	r.mutexFinished.Lock()
	defer r.mutexFinished.Unlock()
	r.isFinished[t.Name()] = true
}

func (r *Scheduler) Interrupt() {
	r.cancel()
}

func (r *Scheduler) Release() ([]types.ITask, error) {
	r.wg.Wait()
	return r.finishedList, r.err
}

func (r *Scheduler) Run(ctx context.Context) error {
	untag := r.AddTag("RunSync")
	r.wg.Add(1)
	defer func() {
		r.wg.Done()
		untag()
	}()
	if r.numOfConcurrencies == 0 {
		r.err = r.startSyncRun(ctx, r.Tasks)
		return r.err
	}

	go r.startProducer()
	go r.startConsumer(ctx)

	r.Log("Waiting for all tasks to finish")
	select {
	case <-r.ctx.Done():
		r.err = r.ctx.Err()
	case err := <-r.chErr:
		r.err = err
	}
	return r.err
}

func (r *Scheduler) startSyncRun(ctx context.Context, tasks []types.ITask) error {
	for _, task := range tasks {
		if r.ctx.Err() != nil {
			return r.ctx.Err()
		}
		key := task.Name()
		logKey := utils.CompactName(key)
		if isFinished, ok := r.isFinished[key]; !ok || !isFinished {
			r.Log("Inspect to run service %s which has %d dependencies", logKey, len(task.DependOn(r.status)))
			if err := r.startSyncRun(ctx, task.DependOn(r.status)); err != nil {
				r.LogS("Service %s failed to run dependencies: %s", logKey, err.Error())
				return err
			}

			if err := task.Run(ctx, r.status); utils.WrapCommonError(err) != nil {
				r.LogS("Service %s failed to %s: %s", logKey, r.status.String(), err.Error())
				return err
			}
			r.LogS("Service %s %s successfully", logKey, r.status.String())
			r.isFinished[key] = true
			r.finishedList = append(r.finishedList, task)
		} else {
			r.Log("Service %s has been already finished. Skip!", logKey)
		}
	}
	return nil
}

func (r *Scheduler) startProducer() {
	log := r.Logger.Clone()
	untag := log.AddTag("startProducer")
	defer func() {
		log.Log("Exit producer")
		close(r.chReqSync)
		close(r.chReqAsync)
		untag()
	}()

	r.checkAndLoad(r.Tasks)

	utils.WaitOnEvents(r.ctx, func(ctx context.Context, task types.ITask) error {
		key := task.Name()
		untag := log.AddTag("response-" + utils.CompactName(key))
		defer untag()
		r.finishedList = append(r.finishedList, task)
		r.mutexFinished.Lock()
		r.isFinished[key] = true
		r.mutexFinished.Unlock()
		if len(r.finishedList) == len(r.Tasks) {
			return utils.ErrorEndOfProcessing
		}
		followers := task.Followers(r.status)
		r.checkAndLoad(followers)
		return nil
	}, nil, r.chRes)
}

func (r *Scheduler) startConsumer(ctx context.Context) {
	log := r.Logger.Clone()
	untag := log.AddTag("startConsumer")
	defer func() {
		log.Log("Exit consumer")
		close(r.chRes)
		close(r.chErr)
		untag()
	}()
	var wgTask sync.WaitGroup
	wgTask.Add(1)
	go r.startSyncWorker(ctx, &wgTask)
	wgTask.Add(1)
	go r.startAsyncWorker(ctx, &wgTask)
	wgTask.Wait()
}

func (r *Scheduler) startSyncWorker(ctx context.Context, wg *sync.WaitGroup) {
	log := r.Logger.Clone()
	untag := log.AddTag("startSyncWorker")
	defer func() {
		log.Log("Exit sync worker")
		wg.Done()
		untag()
	}()
	utils.WaitOnEvents(r.ctx, func(_ context.Context, task types.ITask) error {
		key := task.Name()
		logKey := utils.CompactName(key)
		log.Log("Service %s is going to be run at %s mode", logKey, r.status.String())
		if err := task.Run(ctx, r.status); utils.WrapCommonError(err) != nil {
			log.LogS("Service %s failed to %s: %s", logKey, r.status.String(), err.Error())
			r.chErr <- err
			return err
		}
		log.LogS("Service %s %s successfully", logKey, r.status.String())
		r.chRes <- task
		return nil
	}, r.chErr, r.chReqSync)
}

func (r *Scheduler) startAsyncWorker(ctx context.Context, wg *sync.WaitGroup) {
	log := r.Logger.Clone()
	untag := log.AddTag("startAsyncWorker")
	defer func() {
		log.Log("Exit async worker")
		wg.Done()
		untag()
	}()
	utils.WaitOnEvents(r.ctx, func(_ context.Context, task types.ITask) error {
		wg.Add(1)
		go func(task types.ITask) {
			defer wg.Done()
			key := task.Name()
			logKey := utils.CompactName(key)
			log.Log("Service %s is going to be run at %s mode", logKey, r.status.String())
			if err := task.Run(ctx, r.status); utils.WrapCommonError(err) != nil {
				log.LogS("Service %s failed to %s: %s", logKey, r.status.String(), err.Error())
				r.chErr <- err
				return
			}
			log.LogS("Service %s %s successfully", logKey, r.status.String())
			r.chRes <- task
		}(task)
		return nil
	}, r.chErr, r.chReqAsync)
}

func (r *Scheduler) checkAndLoad(tasks []types.ITask) {
	for _, task := range tasks {
		if r.checkDependenciesReady(task) {
			logKey := utils.CompactName(task.Name())
			r.Log("Service %s is ready to run", logKey)
			key := task.Name()
			r.mutexRun.RLock()
			if isRunning, ok := r.isRunning[key]; !ok || !isRunning {
				r.mutexRun.RUnlock()
				r.mutexRun.Lock()
				r.isRunning[key] = true
				r.mutexRun.Unlock()
				r.ranList = append(r.ranList, task)
				if task.IsRunAsync(r.status) {
					r.Log("Push service %s to async queue", logKey)
					r.chReqAsync <- task
				} else {
					r.Log("Push service %s to sync queue", logKey)
					r.chReqSync <- task
				}
			} else {
				r.mutexRun.RUnlock()
			}
		}
	}
}

func (r *Scheduler) checkDependenciesReady(task types.ITask) bool {
	r.mutexFinished.RLock()
	defer r.mutexFinished.RUnlock()
	logKey := utils.CompactName(task.Name())
	r.Log("Check if service %s is ready to run", logKey)
	for _, dep := range task.DependOn(r.status) {
		depKey := dep.Name()
		if isFinished, ok := r.isFinished[depKey]; !ok || !isFinished {
			r.Log("Service %s is not ready to run because service %s is not ready", logKey, utils.CompactName(depKey))
			return false
		}
	}
	return true
}
