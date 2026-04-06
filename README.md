# PipePool

PipePool is Phase Beta.

The project is intentionally boring on purpose:

- take text input from memory or files
- normalize it
- mark it valid or invalid
- send it into a bounded worker pool
- process it
- collect one final summary

The point is not product design.
The point is learning how one small concurrent system fits together.

## The system picture

```text
main/app
  |
  v
ingest -> normalize -> enqueue -> [ bounded jobs channel ] -> fixed worker pool -> results -> summary collector
```

There are only two concurrent shapes you are practicing here:

1. a small pipeline
2. a bounded worker pool

The summary collector should stay simple. It can run in the main flow while ranging over results.

## How to use these docs

Read them in this order:

1. `docs/NEW_AFTER_ALPHA.md`
2. `docs/FILE_MAP.md`
3. `docs/PHASE_BETA_PIPEPOOL.md`

## What "success" looks like

You are done when all of these feel normal instead of scary:

- one root `context.Context`
- a fixed worker count
- a bounded queue
- visible backpressure
- clean cancellation
- no goroutine leaks
- logs that explain lifecycle clearly

## Ground rule for yourself

Do not try to make PipePool clever.

If you are unsure, choose the simpler version:

- in-memory inputs before files
- one `Config` struct
- one jobs channel
- one results channel
- one summary struct
- one logger

That is enough.
