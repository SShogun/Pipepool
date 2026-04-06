# PipePool

PipePool is Phase Beta.

This project takes the small concurrency shapes from Alpha and puts them into one tiny system:

- ingest text
- normalize and validate it
- push it through a bounded jobs queue
- process it with a fixed worker pool
- collect one final summary

The work stays intentionally boring on purpose.
The point is learning backpressure, ownership, cancellation, and clean shutdown without extra product noise.

## Read This Folder In Order

1. `docs/NEW_AFTER_ALPHA.md`
2. `docs/FILE_MAP.md`
3. `docs/PHASE_BETA_PIPEPOOL.md`

## What Success Looks Like

- one root `context.Context`
- one bounded queue
- one fixed worker count
- clear logs
- visible backpressure
- no goroutine leaks

If Alpha taught clean wires, PipePool is where those wires become one small concurrent program.
