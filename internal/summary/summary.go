package summary

import (
	. "Pipepool/internal/types"
	"context"
)

func Collect(ctx context.Context, results <-chan Result) Summary {
	summary := Summary{}

	for {
		select {
		case <-ctx.Done():
			summary.Errors = append(summary.Errors, ctx.Err())
			return summary
		case result, ok := <-results:
			if !ok {
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
