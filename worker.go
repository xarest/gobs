package gobs

import "context"

type status int

const (
	workerIdle status = iota
	workerBusy
	workerWaiting
)

type WorkerCtl struct {
	InQueue  chan *Component
	ErrQueue chan error
}

func createWorker(
	scheduleCtx context.Context,
	taskCtx context.Context,
	task func(ctx context.Context, c *Component, onError func(error)),
	qBufferSize, limitThread int,
) *WorkerCtl {
	inQueue := make(chan *Component, qBufferSize)
	errQueue := make(chan error, qBufferSize)
	workerController := WorkerCtl{
		InQueue:  inQueue,
		ErrQueue: errQueue,
	}

	for i := 0; i < limitThread; i++ {
		go func(workerId int) {
			for {
				select {
				case <-scheduleCtx.Done():
					return
				case c, ok := <-inQueue:
					if !ok {
						break
					}
					task(taskCtx, c, func(err error) {
						if scheduleCtx.Err() == nil {
							select {
							case <-scheduleCtx.Done():
							case errQueue <- err:
							}
						}
					})
				}
			}
		}(i)
	}
	return &workerController
}

func (w *WorkerCtl) Close() {
	close(w.ErrQueue)
	close(w.InQueue)
}
