package main

import (
	"context"
	"github.com/foojank/fzz"
	"os"
)

func main() {
	switch os.Args[1] {
	case "build":
		err := fzz.Build(os.Args[2], fzz.BuildOptions{
			Output: "fizzpit.fzz",
		})
		if err != nil {
			panic(err)
		}
	case "exec":
		err := fzz.Exec(context.TODO(), os.Args[2], fzz.ExecOptions{
			Command: os.Args[3],
		})
		if err != nil {
			panic(err)
		}
	}
}
