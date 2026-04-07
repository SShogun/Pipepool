package pipeline

import (
	. "Pipepool/internal/types"
)

type Item struct {
	ID     int
	Input  string
	Valid  bool
	Result Result
}
