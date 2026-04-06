# How To Build Phase Beta: PipePool

This guide is written for "I know Alpha, but I do not yet feel comfortable after Alpha."

So the plan is:

- build in small phases
- check your understanding after each phase
- only take hints when you are stuck

Do not try to finish the whole project in one sitting.

## The boring domain to choose

Choose the simplest possible work item:

- input text
- normalized text
- valid or invalid flag
- word count
- line count
- processing duration

That is enough for PipePool.

## The target system

```text
source inputs
  -> ingest
  -> normalize
  -> enqueue into bounded jobs channel
  -> fixed worker pool processes jobs
  -> results channel
  -> summary collector
```

## Phase 0: Set up the repo

Goal:
Make the project buildable before it is interesting.

Create:

- `go.mod`
- `cmd/pipepool/main.go`
- `internal/app/`
- `internal/logging/`
- `internal/pipeline/`
- `internal/pool/`
- `internal/summary/`
- `internal/testutil/`

Self-check:

- can `go test ./...` run without import errors?

Hint if stuck:
Start with tiny packages and empty structs if needed. You are only creating the shape here.

## Phase 1: Define your core types

Goal:
Make the data flow visible before making it concurrent.

You need a few simple structs:

- raw input
- normalized job
- worker result
- final summary
- config

Possible mental split:

- pipeline types describe work before processing
- pool result describes processed output
- summary type describes the end report

Self-check:

- can you explain what changes between raw input and job?
- can you explain what changes between job and result?

Hint if stuck:
If a field is only useful after processing, it belongs on `Result`, not `Job`.

## Phase 2: Build the pipeline without worrying about the pool yet

Goal:
Practice the first concurrent shape by itself.

Implement three stages:

1. `Ingest`
2. `Normalize`
3. `Enqueue`

What each one should do:

- `Ingest` emits raw items from a slice of strings or file contents
- `Normalize` trims whitespace, standardizes line endings, and marks valid/invalid
- `Enqueue` sends prepared jobs into the bounded jobs channel

Important:
Do not mix all three responsibilities into one goroutine if you want to learn the stage boundaries clearly.

Self-check:

- does each stage read from one input and write to one output?
- does each stage stop on `ctx.Done()`?
- does each stage close only the channel it created?

Hint if stuck:
For each stage, ask two questions:

1. What channel does this stage read from?
2. What channel does this stage create and return?

The returned one is the one it should close.

## Phase 3: Add the bounded jobs channel

Goal:
Create real backpressure.

This is the main handoff:

```go
jobs := make(chan Job, cfg.QueueSize)
```

Why it matters:
This is what makes slow workers affect upstream behavior.

How to see it:

- use a very small queue size like `1`
- make workers sleep briefly
- send several inputs
- log when enqueue starts and when it actually succeeds

Self-check:

- if workers are slow, can enqueue visibly block?
- if queue size is large, does the pipeline run ahead more easily?

Hint if stuck:
If backpressure is invisible, your workers are probably too fast or your queue is too big.

## Phase 4: Build the fixed worker pool

Goal:
Practice the second concurrent shape by itself.

Rules:

- start exactly `cfg.WorkerCount` workers
- each worker loops over the jobs channel
- each worker emits one result per job
- when all workers stop, close the results channel

The processing work should stay boring:

- maybe sleep for a configured delay
- count words
- count lines
- return duration
- preserve validity info from the job

Self-check:

- can you point to the exact line where workers are started?
- can you count the number of worker goroutines from the code without guessing?
- is there any hidden `go func()` inside processing?

Hint if stuck:
The pool usually needs a `sync.WaitGroup`.
Workers call `Done`.
A separate goroutine waits, then closes `results`.

## Phase 5: Collect results in one place

Goal:
Finish the flow with one simple collector.

Your collector should probably:

- range over `results`
- count completed jobs
- count invalid jobs
- total words
- total lines
- track slowest job
- gather errors

Important:
Keep this collector simple.
It does not need its own goroutine unless you have a very specific reason.

Self-check:

- can `main` or `app.Run` call one function like `summary.Collect(ctx, results)`?

Hint if stuck:
If you are starting extra goroutines for the collector, you are probably making it harder than needed.

## Phase 6: Add structured logs

Goal:
Make the lifecycle readable without a debugger.

Use fields like:

- `component`
- `job_id`
- `state`
- `duration`
- `worker_id`
- `queue_size`

Important log events:

- app start
- pipeline start
- normalize done
- enqueue waiting
- enqueue success
- worker start
- worker got job
- worker finished job
- cancel received
- results closed
- app stop

Self-check:

- can you read logs and tell why the system stopped?
- can you tell whether a job was slow because of queue wait or processing time?

Hint if stuck:
Use `slog` with `InfoContext` and consistent field names. Consistency matters more than clever log messages.

## Phase 7: Add cancellation

Goal:
Make the whole system stop cleanly.

Use:

- one root context in `main`
- maybe `context.WithTimeout`
- maybe per-job timeout inside worker processing

Where to check `ctx.Done()`:

- inside pipeline sends
- inside worker receive loops
- inside job processing if processing can block
- inside result sends

Self-check:

- if the root context times out, do pipeline stages stop?
- do workers stop?
- does the results channel close?
- does the collector return?

Hint if stuck:
Any blocking send or receive without a `select` around `ctx.Done()` is a place to inspect.

## Phase 8: Make backpressure visible and testable

Goal:
Prove that the queue is not just decorative.

Easy test shape:

- queue size `1`
- worker count `1`
- processing sleeps long enough to create pressure
- submit several jobs quickly

Things to verify:

- enqueue does not race ahead forever
- total runtime reflects queue waiting
- logs show when enqueue had to wait

Hint if stuck:
Backpressure is easiest to see when the system is deliberately tiny and slow.

## Phase 9: Prove clean shutdown

Goal:
Check the "done means done" part of the system.

You want tests for:

- `go test -race ./...`
- worker count never exceeds config
- cancellation returns promptly
- goroutines return near baseline after shutdown

A very practical pattern:

1. record baseline goroutine count
2. run the app
3. cancel or let it finish
4. wait a small moment
5. compare goroutine count again

Hint if stuck:
If goroutines stay alive, suspect:

- a sender waiting forever
- a receiver waiting forever
- a channel never being closed

## The simplest successful version

If you want the smallest version that still teaches the lesson, make this:

- input is `[]string`
- normalize with `TrimSpace` and line-ending cleanup
- invalid means empty after trimming
- workers count words and lines
- summary tracks totals and slowest job
- logs are text logs through `slog`

That version is absolutely enough.

## Hints for wiring without giving away the full answer

### If you are confused about package boundaries

Ask:
Who owns the lifetime of this thing?

- top-level lifetime: `app`
- job preparation lifetime: `pipeline`
- worker lifetime: `pool`
- final report lifetime: `summary`

### If you are confused about channel direction

Ask:
Who sends, who receives, who closes?

That usually gives you the function signature too:

- sender-only argument: `chan<- T`
- receiver-only argument: `<-chan T`

### If you are confused about where the collector runs

Keep it in the main flow first.
Do not add concurrency just because the rest of the program uses concurrency.

### If you are confused about logs

Do not try to log everything.
Log state transitions:

- started
- waiting
- received
- processing
- done
- canceled
- stopped

## Suggested reading order while building

Read only what helps the current phase.

Before Phase 0:

- <https://go.dev/doc/modules/layout>

Before Phase 2:

- <https://go.dev/blog/pipelines>

Before Phase 6:

- <https://pkg.go.dev/log/slog>

Before Phase 7:

- <https://pkg.go.dev/context>

Before Phase 9:

- <https://go.dev/doc/articles/race_detector.html>

## Final reminder

The main win of PipePool is not "I built a text processor."

The main win is:

"I can point to every goroutine, every channel owner, every shutdown path, and every concurrency boundary."

That is the actual Phase Beta skill.
