package pool

import (
	. "Pipepool/internal/types"
	"context"
)

func Run(ctx context.Context, jobs <-chan Job, results chan<- Result) {

}
