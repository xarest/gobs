package gobs

import (
	"context"
	"errors"
	"runtime"

	"github.com/traphamxuan/gobs/logger"
)

type Scheduler struct {
	*logger.Logger
	concurrency int
	services    []*Component
	status      ServiceStatus

	ctx     context.Context
	cancel  context.CancelFunc
	chSync  chan *Component
	chAsync chan *Component
	waiter  chan struct{}
}

func NewScheduler(concurrency int, log *logger.Logger) *Scheduler {
	return &Scheduler{
		Logger:      log,
		concurrency: concurrency,
		status:      StatusInit,
	}
}

func scan(services []*Component, ss ServiceStatus) (sync []*Component, async []*Component) {
	for _, service := range services {
		if err := service.ShouldRun(ss); err != nil {
			continue
		}
		if service.IsRunAsync(ss) {
			async = append(async, service)
		} else {
			sync = append(sync, service)
		}
	}
	return sync, async
}

func (s *Scheduler) Load(services []*Component) {
	s.status = StatusInit
	s.services = services
	if s.chSync != nil {
		close(s.chSync)
	}
	if s.chAsync != nil {
		close(s.chAsync)
	}
	s.chSync = make(chan *Component, len(services))
	s.chAsync = make(chan *Component, len(services))
}

func (s *Scheduler) Run(ctx context.Context, ss ServiceStatus) error {
	untag := s.AddTag("Run")
	defer untag()
	s.ctx, s.cancel = context.WithCancel(ctx)
	if s.concurrency == 0 {
		return s.RunSync(ctx, ss)
	}
	return s.RunAsync(ctx, ss, s.concurrency)
}

func (s *Scheduler) Stop(ctx context.Context) {
	untag := s.AddTag("Stop")
	defer untag()
	if s.cancel != nil && s.ctx.Err() == nil {
		s.cancel()
	}
}

func (s *Scheduler) RunSync(ctx context.Context, ss ServiceStatus) error {
	untag := s.AddTag("RunSync")
	defer untag()
	if s.status >= ss {
		return nil
	}
	s.status = ss
	sync, async := scan(s.services, ss)
	services := append(sync, async...)
	for _, service := range services {
		if err := s.executeAll(ctx, service); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) RunAsync(ctx context.Context, ss ServiceStatus, concurrence int) (err error) {
	untag := s.AddTag("RunAsync")
	// s.SetDetail(true)
	if s.status >= ss {
		return nil
	}
	totalServices := s.countServiceRun(ss)
	if concurrence <= 0 {
		concurrence = totalServices
		if concurrence > runtime.NumCPU() {
			concurrence = runtime.NumCPU()
		}
	}
	s.status = ss
	s.waiter = make(chan struct{}, concurrence)
	for i := 0; i < concurrence; i++ {
		s.waiter <- struct{}{}
	}

	processor := func(ctx context.Context, service *Component, onError func(error)) {
		log := s.Logger.Clone()
		untag := log.AddTag(compactName(service.name))
		defer func() {
			untag()
		}()
		if err := s.executeParallel(ctx, service); err != nil {
			onError(err)
			return
		}
	}

	syncWorker := createWorker(s.ctx, ctx, processor, len(s.services), 1)
	s.chSync = syncWorker.InQueue

	asyncWorker := createWorker(s.ctx, ctx, processor, len(s.services), concurrence)
	s.chAsync = asyncWorker.InQueue

	defer func() {
		s.cancel()
		syncWorker.Close()
		asyncWorker.Close()
		close(s.waiter)
		s.chSync = nil
		s.chAsync = nil
		s.waiter = nil
		untag()
	}()

	sync, async := scan(s.services, s.status)
	for _, service := range sync {
		select {
		case <-s.ctx.Done():
			return
		case s.chSync <- service:
		}
	}
	for _, service := range async {
		select {
		case <-s.ctx.Done():
			return
		case s.chAsync <- service:
		}
	}

	for wrapCommonError(err) == nil {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case err = <-syncWorker.ErrQueue:
		case err = <-asyncWorker.ErrQueue:
		case s.waiter <- struct{}{}:
			totalServices--
			if totalServices == 0 {
				return nil
			}
		}
	}
	return err
}

func (s *Scheduler) executeParallel(ctx context.Context, service *Component) (err error) {
	// untag := s.AddTag("executeParallel")
	// defer untag()
	onSync := func(ctx context.Context, dep *Component) error {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case s.chSync <- dep:
		}
		return nil
	}
	onAsync := func(ctx context.Context, dep *Component) error {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case s.chAsync <- dep:
		}
		return nil
	}
	return s.execute(ctx, service, &onSync, &onAsync)
}

func (s *Scheduler) executeAll(ctx context.Context, service *Component) (err error) {
	untag := s.AddTag("executeAll")
	defer untag()
	process := func(ctx context.Context, dep *Component) error {
		return s.executeAll(ctx, dep)
	}
	return s.execute(ctx, service, &process, &process)
}

func (s *Scheduler) execute(
	ctx context.Context,
	service *Component,
	onSyncFollower *func(context.Context, *Component) error,
	onAsyncFollower *func(context.Context, *Component) error,
) (err error) {
	// untag := s.AddTag("execute")
	// defer untag()
	err = service.ShouldRun(s.status)
	if err == nil {
		s.Log("Running service %s at status %s", service.name, s.status.String())
		if err = service.Run(ctx, s.status); s.logRun(service, err) != nil {
			return err
		}
		if err == nil && s.waiter != nil {
			s.Log("Trying to notify about service %s", service.name)
			select {
			case <-s.ctx.Done():
				return s.ctx.Err()
			case <-s.waiter:
			}
		}
	} else if errors.Is(err, ErrorServiceNotReady) {
		s.Log("Service %s is not ready to run at status %s", service.name, s.status.String())
		return err
	}

	if onSyncFollower != nil || onAsyncFollower != nil {
		followers := service.Followers(s.status)
		s.Log("Running followers of service [%d]", len(followers))
		sync, async := scan(followers, s.status)
		if onSyncFollower != nil {
			for _, dep := range sync {
				if err := (*onSyncFollower)(ctx, dep); wrapCommonError(err) != nil {
					return err
				}
			}
		}

		if onAsyncFollower != nil {
			for _, dep := range async {
				if err := (*onAsyncFollower)(ctx, dep); wrapCommonError(err) != nil {
					return err
				}
			}
		}
	}
	return err
}

func (s *Scheduler) countServiceRun(ss ServiceStatus) int {
	count := 0
	for _, service := range s.services {
		if service.status >= ss {
			continue
		}
		if service.status < StatusSetup && ss >= StatusStop {
			continue
		}
		count++
	}
	return count
}

func (s *Scheduler) logRun(service *Component, err error) error {
	status := "successfully"
	formatStr := "service %s %s %s"
	var step string
	switch s.status {
	case StatusInit:
		step = "inits"
	case StatusSetup:
		step = "setups"
	case StatusStart:
		step = "starts"
	case StatusStop:
		step = "stops"
	}

	if err == ErrorServiceRan {
		status = "skipped"
		s.LogS(formatStr, service.name, step, status)
		err = nil
	} else if err != nil {
		status = "failed"
		formatStr += ", %v"
		s.LogS(formatStr, service.name, step, status, err)
	} else {
		s.LogS(formatStr, service.name, step, status)
	}

	return err
}
