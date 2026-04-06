# New After Alpha

This file is for things that Phase Alpha probably did not make fully normal yet.

If a word below feels fuzzy, that is okay.
That is exactly why this file exists.

## 1. `cmd/` and `internal/`

What it means:

- `cmd/pipepool/main.go` is the actual runnable program
- `internal/...` holds helper packages used only inside this module

Why you are using them now:
PipePool is your first "small real repo" instead of one exercise file.

What to read:

- Go module layout: <https://go.dev/doc/modules/layout>
- How Go code is organized: <https://go.dev/doc/code.html>

## 2. Passing `context.Context` across packages

What it means:
Long-running functions should accept `ctx context.Context` as their first parameter.

Why this matters:
If cancellation starts in `main`, every long-running stage should be able to see it.

The rule:

```go
func Run(ctx context.Context, cfg Config) error
```

Not this:

```go
type App struct {
    Ctx context.Context
}
```

Why not store it on a struct:
It hides lifetimes and makes cancellation harder to reason about.

What to read:

- `context` package docs: <https://pkg.go.dev/context>
- Canceling in-progress operations: <https://go.dev/doc/database/cancel-operations>

## 3. Channel ownership across packages

What it means:
The code that creates a channel is usually the code that closes it.

Why this matters:
Once you split into packages, it becomes much easier to close a channel from the wrong place.

The PipePool ownership rule:

- pipeline owns `jobs`
- pool owns `results`
- summary owns neither

What to read:

- Pipelines and cancellation: <https://go.dev/blog/pipelines>

## 4. Bounded worker pool as a hard limit

What it means:
If config says `WorkerCount: 4`, your system should never process with 5 or 40 goroutines "just for a moment".

Why this matters:
This is how you learn to control concurrency instead of letting it spread.

The smell to watch for:
Any `go func()` inside per-job processing is suspicious here.

What to read:

- Pipelines and cancellation: <https://go.dev/blog/pipelines>

## 5. Backpressure

What it means:
When downstream is slow, upstream must slow down too.

In this project:

- workers are slow
- jobs channel fills up
- enqueue blocks
- normalize eventually blocks
- ingest eventually blocks

Why this is good:
It keeps memory bounded and makes pressure visible.

What to read:

- Pipelines and cancellation: <https://go.dev/blog/pipelines>

## 6. Structured logging

What it means:
You log key/value fields, not just pretty sentences.

Example:

```go
logger.InfoContext(ctx, "job finished",
    "component", "pool",
    "job_id", job.ID,
    "state", "done",
    "duration", dur,
)
```

Why this matters:
Concurrent systems become much easier to understand when every log line has consistent fields.

What to read:

- `log/slog` docs: <https://pkg.go.dev/log/slog>

## 7. Race-clean testing

What it means:
`go test -race ./...` should pass.

Why this matters:
Concurrency bugs often look fine until the race detector catches shared-state mistakes.

What to read:

- Race detector: <https://go.dev/doc/articles/race_detector.html>

## 8. Goroutine leak checking

What it means:
After shutdown, goroutines should return close to baseline instead of staying alive forever.

Why this matters:
A blocked send, blocked receive, or forgotten goroutine often shows up as a leak.

What to watch for:

- send on a channel with no receiver
- receive from a channel that is never closed
- goroutine waiting forever after cancellation

You may not know this yet:
`runtime.NumGoroutine()` is not a perfect truth machine, but it is a useful safety check in small projects.

## 9. `select` now means lifecycle control, not just syntax practice

In Alpha, `select` may have felt like an isolated feature.

In PipePool, `select` becomes part of the system contract:

- send job or stop on cancel
- receive result or stop on cancel
- wait for work or exit when input closes

This is the point where the pieces start becoming one system.

## What is genuinely new for you here

The new difficulty is not "harder Go syntax".

The new difficulty is:

- package boundaries
- ownership boundaries
- lifetime boundaries
- pressure moving across stages
- logs explaining what happened

That is the real Beta jump.
