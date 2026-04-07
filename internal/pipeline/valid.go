package pipeline

import (
	"context"
)

func validate(ctx context.Context, items <-chan Item) <-chan Item {
	out := make(chan Item)

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

				item.Valid = item.Input != ""

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
