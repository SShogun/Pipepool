package summary

import (
	"Pipepool/internal/types"
	"context"
	"io"
	"log/slog"
)

func Collect(ctx context.Context, results <-chan types.Result, logger *slog.Logger) types.Summary {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	summary := types.Summary{}
	logger.InfoContext(ctx, "summary lifecycle", "component", "summary", "state", "start")

	for {
		select {
		case <-ctx.Done():
			summary.Errors = append(summary.Errors, ctx.Err())
			logger.InfoContext(ctx, "summary lifecycle", "component", "summary", "state", "canceled", "total_jobs", summary.TotalJobs)
			return summary
		case result, ok := <-results:
			if !ok {
				logger.InfoContext(ctx, "summary lifecycle", "component", "summary", "state", "stop", "total_jobs", summary.TotalJobs, "successes", summary.SuccessCount, "failures", summary.FailureCount)
				return summary
			}

			summary.TotalJobs++
			summary.TotalWords += result.WordCount
			summary.TotalLines += result.LineCount
			summary.TotalDuration += result.Duration

			if result.Duration > summary.MaxDuration {
				summary.MaxDuration = result.Duration
				summary.SlowestJobID = result.ID
			}

			if result.Valid && result.Err == nil {
				summary.SuccessCount++
			} else {
				summary.FailureCount++
			}

			if result.Err != nil {
				summary.Errors = append(summary.Errors, result.Err)
			}
		}
	}
}
