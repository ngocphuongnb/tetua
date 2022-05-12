package jadecmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Joker/jade"
)

func reset() {
	dict = make(map[string]string)
	lib_name = ""
	outdir = ""
	basedir = ""
	pkg_name = "views"
	stdlib = false
	stdbuf = false
	writer = false
	inline = false
	format = false
	ns_files = make(map[string]bool)
}

type CompileConfig struct {
	OutDir  string
	BaseDir string
	PkgName string
	Stdlib  bool
	Stdbuf  bool
	Writer  bool
	Inline  bool
	Format  bool
}

func CompilePaths(jadePaths []string, config CompileConfig) {
	reset()
	outdir = config.OutDir
	basedir = config.BaseDir
	pkg_name = config.PkgName
	stdlib = config.Stdlib
	stdbuf = config.Stdbuf
	writer = config.Writer
	inline = config.Inline
	format = config.Format

	jade.Config(golang)

	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.MkdirAll(outdir, 0755)
	}
	outdir, _ = filepath.Abs(outdir)

	if _, err := os.Stat(basedir); !os.IsNotExist(err) && basedir != "./" {
		os.Chdir(basedir)
	}

	for _, jadePath := range jadePaths {
		stat, err := os.Stat(jadePath)
		if err != nil {
			log.Fatalln(err)
		}

		absPath, _ := filepath.Abs(jadePath)
		if stat.IsDir() {
			genDir(absPath, outdir, pkg_name)
		} else {
			genFile(absPath, outdir, pkg_name)
		}
		if !stdlib {
			makeJfile(stdbuf)
		}
	}
}
