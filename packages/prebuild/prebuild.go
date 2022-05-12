package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/ngocphuongnb/tetua/app/config"
	jadecmd "github.com/ngocphuongnb/tetua/packages/jade"
)

type BuildCache struct {
	ViewMTime map[string]int64
	EntMTime  int64
}

func getDirLastMTime(root string) int64 {
	var lastModTime int64
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.ModTime().UnixMicro() > lastModTime {
				lastModTime = info.ModTime().UnixMicro()
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return lastModTime
}

func getBuildCache(buildCacheFile string) BuildCache {
	buildCache := BuildCache{}
	buildCacheBytes, err := ioutil.ReadFile(buildCacheFile)

	if err == nil {
		if err = json.Unmarshal(buildCacheBytes, &buildCache); err != nil {
			log.Fatal(err)
		}
	}

	return buildCache
}

func BuildViewAssets(viewsDir, viewsOutputDir, buildCacheFile string, force bool) error {
	buildCache := getBuildCache(buildCacheFile)
	cachedViewMTimes := buildCache.ViewMTime
	changedPageFiles := make([]string, 0)
	changedPartialsFiles := make([]string, 0)
	changedFiles := make([]string, 0)

	if !force {
		if err := filepath.Walk(viewsDir, func(viewFilePath string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				isPageFile := strings.HasPrefix(viewFilePath, path.Join(viewsDir, "pages"))
				fileMTime := info.ModTime().UnixMicro()

				if viewFileLastMTime, ok := cachedViewMTimes[viewFilePath]; ok && viewFileLastMTime >= fileMTime {

				} else {
					if isPageFile {
						changedPageFiles = append(changedPageFiles, viewFilePath)
					} else {
						changedPartialsFiles = append(changedPartialsFiles, viewFilePath)
					}
					cachedViewMTimes[viewFilePath] = fileMTime
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// If has partials change, rebuild all
	if force || len(changedPartialsFiles) > 0 {
		fmt.Println("> Has partials changed")
		changedFiles = append(changedFiles, viewsDir+"/pages")
	} else {
		// If has page change, rebuild only changed page
		if len(changedPageFiles) > 0 {
			fmt.Println("> Has pages changed")
			changedFiles = changedPageFiles
		}
	}

	if len(changedFiles) > 0 {
		jadecmd.CompilePaths(changedFiles, jadecmd.CompileConfig{
			Writer:  true,
			PkgName: "views",
			OutDir:  viewsOutputDir,
		})

		buildCache.ViewMTime = cachedViewMTimes
		b, err := json.Marshal(buildCache)
		if err != nil {
			return err
		}

		if err = os.WriteFile(buildCacheFile, b, 0644); err != nil {
			return err
		}
	} else {
		fmt.Println("> No views changed")
	}

	return nil
}

func GenerateEnt(buildCacheFile string, force bool) error {
	lastEntSchemaMTime := getDirLastMTime("./packages/entrepository/ent/schema")
	buildCache := getBuildCache(buildCacheFile)

	if force || lastEntSchemaMTime > buildCache.EntMTime {
		fmt.Println("> Has ent schema changed")
		buildCache.EntMTime = lastEntSchemaMTime
		cmd := exec.Command("go", "generate", "./packages/entrepository/ent", "-vvvv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}

		b, err := json.Marshal(buildCache)
		if err != nil {
			return err
		}

		if err = os.WriteFile(buildCacheFile, b, 0644); err != nil {
			return err
		}
		return nil
	}

	fmt.Println("> No ent schema changed")
	return nil
}

func main() {
	force := false
	if len(os.Args) > 1 {
		force = os.Args[1] == "--force"
	}

	if err := BuildViewAssets(
		path.Join("./app/themes", config.APP_THEME, "views"),
		"./views",
		"./private/tmp/cache.json",
		force,
	); err != nil {
		panic(err)
	}
	if err := GenerateEnt("./private/tmp/cache.json", force); err != nil {
		panic(err)
	}
}
