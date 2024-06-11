package gobs

import "context"

func createWorker(
	ctx context.Context,
	task func(ctx context.Context, c *Component, onFinish func(error)),
	qBufferSize, limitThread int,
) (inQueue, outQueue chan *Component, errQueue chan error) {
	inQueue = make(chan *Component, qBufferSize)
	outQueue = make(chan *Component, qBufferSize)
	errQueue = make(chan error, qBufferSize)
	for i := 0; i < limitThread; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case c, ok := <-inQueue:
					if !ok {
						break
					}
					task(ctx, c, func(err error) {
						if err != nil {
							select {
							case <-ctx.Done():
							case errQueue <- err:
							}
						} else {
							select {
							case <-ctx.Done():
							case outQueue <- c:
							}
						}
					})
				}
			}
		}()
	}
	return inQueue, outQueue, errQueue
}
