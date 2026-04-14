package pipeline

import (
	"Pipepool/internal/testutil"
	"context"
	"testing"
	"time"
)

func TestQueueAppliesBackpressure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan Item)
	out := queue(ctx, in, 1, testutil.NewDiscardLogger())

	senderDone := make(chan struct{})
	go func() {
		defer close(senderDone)
		in <- Item{ID: 1, Input: "first", Valid: true}
		in <- Item{ID: 2, Input: "second", Valid: true}
		in <- Item{ID: 3, Input: "third", Valid: true}
		close(in)
	}()

	select {
	case <-senderDone:
		t.Fatal("sender completed without blocking; expected bounded queue backpressure")
	case <-time.After(100 * time.Millisecond):
	}

	readItem(t, out, 1)

	select {
	case <-senderDone:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("sender stayed blocked after queue made space; expected backpressure release")
	}

	readItem(t, out, 2)
	readItem(t, out, 3)

	if _, ok := <-out; ok {
		t.Fatal("queue output channel should be closed after draining all inputs")
	}
}

func readItem(t *testing.T, out <-chan Item, expectedID int) {
	t.Helper()

	select {
	case item, ok := <-out:
		if !ok {
			t.Fatal("output channel closed unexpectedly")
		}
		if item.ID != expectedID {
			t.Fatalf("got item ID %d, want %d", item.ID, expectedID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for item")
	}
}
