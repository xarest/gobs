package gobs

import "context"

func createWorker(
	ctx context.Context,
	task func(ctx context.Context, c *Component, onFinish func(error)),
	qBufferSize, limitThread int,
) (queue chan *Component, chErr chan error, cancel context.CancelFunc) {
	ctxCancel, cancel := context.WithCancel(ctx)
	queue = make(chan *Component, qBufferSize)
	chErr = make(chan error, qBufferSize)
	chLimit := make(chan struct{}, limitThread)
	go func() {
		defer func() {
			cancel()
			// close(chLimit)
			// close(chErr)
		}()
		for {
			select {
			case <-ctxCancel.Done():
				return
			case c := <-queue:
				select {
				case <-ctxCancel.Done():
					return
				case chLimit <- struct{}{}:
					task(ctxCancel, c, func(err error) {
						if err != nil {
							select {
							case <-ctxCancel.Done():
							case chErr <- err:
							}
						}
						select {
						case <-ctxCancel.Done():
						case <-chLimit:
						}
					})
				}
			}
		}
	}()
	return queue, chErr, cancel
}
