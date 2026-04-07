package pipeline

import (
	"context"
)

func queue(ctx context.Context, items <-chan Item, queueSize int) <-chan Item {
	out := make(chan Item, queueSize)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-items:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case out <- item:
				}
			}
		}
	}()

	return out
}
