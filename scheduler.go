package gobs

import (
	"context"
	"errors"
	"sync"

	"github.com/traphamxuan/gobs/logger"
)

type Scheduler struct {
	*logger.Logger
	concurrency int
	services    []*Component
	status      ServiceStatus
	chSync      chan *Component
	chAsync     chan *Component
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
	if s.concurrency == 0 {
		return s.RunSync(ctx, ss)
	}
	return s.RunAsync(ctx, ss, s.concurrency)
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
	if concurrence <= 0 {
		concurrence = len(s.services)
	}
	s.status = ss
	totalService := len(s.services)
	var wg sync.WaitGroup
	wg.Add(totalService)

	processor := func(ctx context.Context, service *Component, onFinish func(error)) {
		log := s.Logger.Clone()
		untag := log.AddTag(compactName(service.name))
		defer func() {
			untag()
		}()
		log.Log("Start service %s (%d)", service.name, totalService)
		if err := s.executeParallel(ctx, service); err != nil {
			onFinish(err)
			return
		}
		totalService -= 1
		log.LogS("Finished service %s (%d)", service.name, totalService)
		wg.Done()
		onFinish(nil)
	}

	if s.chSync != nil {
		close(s.chSync)
	}
	chSync, chErrSync, cancelSync := createWorker(ctx, processor, len(s.services), 1)
	s.chSync = chSync

	if s.chAsync != nil {
		close(s.chAsync)
	}
	chAsync, chErrAsync, cancelAsync := createWorker(ctx,
		func(ctx context.Context, c *Component, onFinish func(error)) {
			go processor(ctx, c, onFinish)
		},
		len(s.services), concurrence,
	)
	s.chAsync = chAsync

	defer func() {
		untag()
		cancelSync()
		cancelAsync()
	}()

	ctxDone, cancel := context.WithCancel(ctx)
	go func(cancel context.CancelFunc) {
		defer cancel()
		sync, async := scan(s.services, s.status)
		for _, service := range sync {
			s.chSync <- service
		}
		for _, service := range async {
			s.chAsync <- service
		}
		wg.Wait()
	}(cancel)

	for wrapCommonError(err) == nil {
		select {
		case <-ctxDone.Done():
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case err = <-chErrSync:
		case err = <-chErrAsync:
		}
	}
	return err
}

func (s *Scheduler) executeParallel(ctx context.Context, service *Component) (err error) {
	// untag := s.AddTag("executeParallel")
	// defer untag()
	onSync := func(ctx context.Context, dep *Component) error {
		s.chSync <- dep
		return nil
	}
	onAsync := func(ctx context.Context, dep *Component) error {
		s.chAsync <- dep
		return nil
	}
	return s.execute(ctx, service, &onSync, &onAsync)
}

func (s *Scheduler) executeAll(ctx context.Context, service *Component) (err error) {
	untag := s.AddTag("executeAll")
	defer untag()
	onSync := func(ctx context.Context, dep *Component) error {
		return s.executeAll(ctx, dep)
	}
	onAsync := func(ctx context.Context, dep *Component) error {
		return s.executeAll(ctx, dep)
	}
	return s.execute(ctx, service, &onSync, &onAsync)
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
