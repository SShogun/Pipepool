# PipePool File Map

This file answers one question:

What should each file be responsible for?

You do not need to create every file on day one.
Start with the "minimum first pass" layout, then split files later if they grow.

## Minimum first pass

```text
Pipepool/
  go.mod
  cmd/
    pipepool/
      main.go
  internal/
    app/
      app.go
      config.go
    logging/
      logger.go
    pipeline/
      pipeline.go
    pool/
      pool.go
    summary/
      summary.go
    testutil/
      testutil.go
```

## Expanded layout

This is the version you can grow into after the first pass feels stable.

```text
Pipepool/
  go.mod
  README.md
  cmd/
    pipepool/
      main.go
  internal/
    app/
      app.go
      config.go
    logging/
      logger.go
    pipeline/
      types.go
      ingest.go
      normalize.go
      enqueue.go
    pool/
      types.go
      run.go
      process.go
    summary/
      types.go
      collect.go
    testutil/
      fixtures.go
      logcapture.go
```

## File by file

### `go.mod`

Purpose:
Declares the module so Go knows how to build packages inside `Pipepool/`.

You may not know this yet:
The module is the build boundary for the project. Every package under this folder belongs to it unless another `go.mod` appears deeper in the tree.

Do not overthink:
Use a simple local module path. This is not the interesting part of the exercise.

### `cmd/pipepool/main.go`

Purpose:
Own the executable entrypoint.

This file should do only a few things:

- create the root context
- create the config
- create the logger
- call `app.Run`
- print or log the final summary

You may not know this yet:
`cmd/` is just a common Go convention for executable programs. `main.go` should stay thin.

### `internal/app/app.go`

Purpose:
Wire the whole system together.

This is the file that should make the project feel like one small machine:

- call pipeline stages
- start the pool
- call the summary collector
- return the final summary or error

You may not know this yet:
This package is not "business logic". It is the ownership and lifetime wiring package.

### `internal/app/config.go`

Purpose:
Hold one `Config` struct and validate it.

Typical fields:

- `WorkerCount int`
- `QueueSize int`
- `RunTimeout time.Duration`
- `PerJobTimeout time.Duration`

You may not know this yet:
A config struct is just one place to keep knobs together. It is not advanced magic.

### `internal/logging/logger.go`

Purpose:
Create a configured `*slog.Logger`.

This file should decide things like:

- text vs JSON handler
- default fields
- log level

You may not know this yet:
Structured logging means "logs as fields", not just formatted English sentences.

Example mental model:

```text
component=pool job_id=4 state=done duration=220ms
```

### `internal/pipeline/pipeline.go`

Purpose:
Hold the full first-pass pipeline in one place until the flow is clear.

This package is responsible for:

- ingesting raw input
- normalizing text
- deciding valid or invalid
- enqueueing jobs into the bounded jobs channel

You may not know this yet:
The queue is the handoff point between the pipeline and the worker pool.

When you later split this file:

- `types.go` holds structs
- `ingest.go` creates raw items
- `normalize.go` cleans text and marks validity
- `enqueue.go` owns the jobs channel send side

### `internal/pool/pool.go`

Purpose:
Own the fixed worker pool.

This package should:

- start exactly `cfg.WorkerCount` worker goroutines
- read jobs from the jobs channel
- process each job
- emit results
- close the results channel after all workers stop

You may not know this yet:
The worker count must be a hard cap. If you start extra goroutines per job, you broke the exercise.

### `internal/summary/summary.go`

Purpose:
Define the result summary and collect results into it.

This package should answer:

- how many jobs were valid
- how many were invalid
- how many words or lines were processed
- which job was slowest
- which errors happened

You may not know this yet:
Your collector does not need to be concurrent just because the earlier parts are concurrent.

### `internal/testutil/testutil.go`

Purpose:
Hold reusable test helpers.

Good uses:

- fake input builders
- slow-job helpers
- log capture helpers
- tiny timeouts for tests

You may not know this yet:
`testutil` exists to keep test code readable, not to look impressive.

## Channel ownership map

This is one of the most important parts of the whole project.

- pipeline creates the jobs channel
- pipeline closes the jobs channel
- pool only receives from the jobs channel
- pool creates the results channel
- pool closes the results channel
- summary only receives from the results channel

If you forget who closes what, come back to this section.

## If the file count feels overwhelming

That feeling is normal.

Start with these packages only:

- `app`
- `logging`
- `pipeline`
- `pool`
- `summary`

Inside each package, begin with one `.go` file.
Split files only after you can explain the flow out loud.
