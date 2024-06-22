package fzz

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

type Imports = interp.Exports

var (
	Stdlib       = stdlib.Symbols
	Syscall      = syscall.Symbols
	Unrestricted = unrestricted.Symbols
	Unsafe       = unsafe.Symbols
)
