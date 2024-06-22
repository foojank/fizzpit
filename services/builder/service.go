package builder

import (
	"archive/zip"
	"context"
	"github.com/otiai10/copy"
	"golang.org/x/mod/modfile"
	"os"
	"os/exec"
	"path/filepath"
)

type Arguments struct {
	Output string
}

type Service struct {
	input string
	args  Arguments
}

func New(input string, args Arguments) *Service {
	return &Service{
		input: input,
		args:  args,
	}
}

func (s *Service) Start(_ context.Context) error {
	// TODO: validate inputs!
	// 	- opts.Output
	// 	- src must be a directory
	tmpDir, err := os.MkdirTemp(".", "fzzpt*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = createDirEnv(tmpDir, s.input)
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
	err = os.Rename(f.Name(), s.args.Output)
	if err != nil {
		return err
	}

	return nil
}

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
