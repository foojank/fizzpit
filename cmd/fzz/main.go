package main

import (
	"context"
	"github.com/foojank/fzz"
	"github.com/foojank/fzz/services/builder"
	"github.com/foojank/fzz/services/executor"
	"os"
)

func main() {
	switch os.Args[1] {
	case "build":
		err := fzz.Build(context.TODO(), os.Args[2], builder.Arguments{
			Output: "fizzpit.fzz",
		})
		if err != nil {
			panic(err)
		}
	case "exec":
		b, err := os.ReadFile(os.Args[2])
		if err != nil {
			panic(err)
		}

		err = fzz.Exec(context.TODO(), b, executor.Arguments{
			Command:      os.Args[3],
			Unrestricted: false,
			Imports: []fzz.Imports{
				fzz.Stdlib,
			},
		})
		if err != nil {
			panic(err)
		}
	}
}
