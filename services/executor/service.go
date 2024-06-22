package executor

import (
	"archive/zip"
	"bytes"
	"context"
	"github.com/traefik/yaegi/interp"
	"io"
)

type Arguments struct {
	Command      string
	Args         []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	Env          []string
	Unrestricted bool
	Imports      []interp.Exports
}

type Service struct {
	r    *bytes.Reader
	args Arguments
}

func New(b []byte, args Arguments) *Service {
	return &Service{
		r:    bytes.NewReader(b),
		args: args,
	}
}

func (s *Service) Start(ctx context.Context) error {
	zr, err := zip.NewReader(s.r, s.r.Size())
	if err != nil {
		return err
	}

	yi := interp.New(interp.Options{
		GoPath:               ".",
		Stdin:                s.args.Stdin,
		Stdout:               s.args.Stdout,
		Stderr:               s.args.Stderr,
		Args:                 s.args.Args,
		Env:                  s.args.Env,
		SourcecodeFilesystem: zr,
		Unrestricted:         s.args.Unrestricted,
	})

	for i := range s.args.Imports {
		_ = yi.Use(s.args.Imports[i])
	}

	_, err = yi.EvalPathWithContext(ctx, s.args.Command)
	return err
}
