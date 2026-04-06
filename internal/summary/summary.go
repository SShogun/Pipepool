package summary

import (
	. "Pipepool/internal/types"
	"context"
)

func Collect(ctx context.Context, results <-chan Result) Summary {
	return Summary{}
}
