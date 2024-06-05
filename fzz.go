package fzz

import (
	"archive/zip"
	"context"
	"github.com/otiai10/copy"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func getModuleName(file string) (string, error) {
	goModBytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	modName := modfile.ModulePath(goModBytes)
	return modName, nil
}

func createDirEnv(dst, src string) error {
	goModPath := filepath.Join(src, "go.mod")
	moduleName, err := getModuleName(goModPath)
	if err != nil {
		return err
	}

	moduleRoot := filepath.Join(dst, "src", moduleName)
	err = os.MkdirAll(moduleRoot, 0755)
	if err != nil {
		return err
	}

	err = copy.Copy(src, moduleRoot)
	if err != nil {
		return err
	}

	return runGoModVendor(moduleRoot)
}

func runGoModVendor(dst string) error {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = dst
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type BuildOptions struct {
	Output string
}

func Build(src string, opts BuildOptions) error {
	// TODO: validate inputs!
	// 	- opts.Output
	// 	- src must be a directory
	tmpDir, err := os.MkdirTemp(".", "fzzpt*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = createDirEnv(tmpDir, src)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp(".", "fzzpt*.fzz")
	if err != nil {
		return err
	}

	zw := zip.NewWriter(f)
	defer os.Remove(f.Name())
	err = zw.AddFS(os.DirFS(tmpDir))
	if err != nil {
		return err
	}

	_ = zw.Close()
	err = os.Rename(f.Name(), opts.Output)
	if err != nil {
		return err
	}

	return nil
}

type ExecOptions struct {
	Command      string
	Args         []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	Env          []string
	Unrestricted bool
}

func Exec(ctx context.Context, file string, opts ExecOptions) error {
	zr, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	defer zr.Close()

	yi := interp.New(interp.Options{
		GoPath:               ".",
		Stdin:                opts.Stdin,
		Stdout:               opts.Stdout,
		Stderr:               opts.Stderr,
		Args:                 opts.Args,
		Env:                  opts.Env,
		SourcecodeFilesystem: zr,
		Unrestricted:         opts.Unrestricted,
	})
	err = yi.Use(stdlib.Symbols)
	if err != nil {
		return err
	}

	if opts.Unrestricted {
		err = yi.Use(unrestricted.Symbols)
		if err != nil {
			return err
		}
	}

	_, err = yi.EvalPathWithContext(ctx, opts.Command)
	if err != nil {
		return err
	}

	return nil
}
