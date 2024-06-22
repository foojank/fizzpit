package fzz

import (
	"context"
	"github.com/foojank/fzz/services/builder"
	"github.com/foojank/fzz/services/executor"
)

func Build(ctx context.Context, src string, args builder.Arguments) error {
	s := builder.New(src, args)
	return s.Start(ctx)
}

func Exec(ctx context.Context, b []byte, args executor.Arguments) error {
	s := executor.New(b, args)
	return s.Start(ctx)
}
