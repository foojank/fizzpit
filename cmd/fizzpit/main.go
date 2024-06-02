package main

import (
	"context"
	"github.com/foojank/fizzpit"
	"os"
)

func main() {
	switch os.Args[1] {
	case "build":
		err := fizzpit.Build(os.Args[2], fizzpit.BuildOptions{
			Output: "fizzpit.fzz",
		})
		if err != nil {
			panic(err)
		}
	case "exec":
		err := fizzpit.Exec(context.TODO(), os.Args[2], fizzpit.ExecOptions{
			Command: os.Args[3],
		})
		if err != nil {
			panic(err)
		}
	}
}
