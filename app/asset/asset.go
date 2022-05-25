package asset

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type StaticAsset struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	DisableMinify bool   `json:"disable_minify"`
	DisableInline bool   `json:"disable_inline"`
}

type ThemeConfig struct {
	Asset struct {
		Dir           string   `json:"dir"`
		DisableMinify []string `json:"disable_minify"`
		DisableInline []string `json:"disable_inline"`
	} `json:"asset"`
}

var themeAssetsBundled = false
var bundledAssets = []*StaticAsset{}
var assets []*StaticAsset

func Load(themeDir string, force bool) {
	loadThemeAssets(themeDir, force)

	for _, bundledAsset := range bundledAssets {
		hasBundledAsset := false
		for _, asset := range assets {
			if bundledAsset.Name == asset.Name {
				asset.Content = bundledAsset.Content
				hasBundledAsset = true
				break
			}
		}

		if !hasBundledAsset {
			assets = append(assets, bundledAsset)
		}
	}
}

func loadThemeAssets(themeDir string, force bool) {
	if themeAssetsBundled && !force && !config.DEVELOPMENT {
		return
	}

	assets = []*StaticAsset{
		{
			Type: "css",
			Name: "editor/tippy-6.3.7.min.css",
			Path: path.Join(config.ROOT_DIR, "packages/editor/dist/tippy-6.3.7.min.css"),
		},
		{
			Type: "css",
			Name: "editor/tippy-light-6.3.7.min.css",
			Path: path.Join(config.ROOT_DIR, "packages/editor/dist/tippy-light-6.3.7.min.css"),
		},
		{
			Type: "css",
			Name: "editor/style.css",
			Path: path.Join(config.ROOT_DIR, "packages/editor/dist/style.css"),
		},
		{
			Type:          "js",
			Name:          "editor/highlight-11.5.0.min.js",
			Path:          path.Join(config.ROOT_DIR, "packages/editor/dist/highlight-11.5.0.min.js"),
			DisableMinify: true,
			DisableInline: true,
		},
		{
			Type:          "js",
			Name:          "editor/editor.js",
			Path:          path.Join(config.ROOT_DIR, "packages/editor/dist/editor.js"),
			DisableMinify: true,
			DisableInline: true,
		},
	}

	themeConfigFile := path.Join(themeDir, "theme.json")
	themeConfig := &ThemeConfig{}
	file, err := ioutil.ReadFile(themeConfigFile)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(file), themeConfig); err != nil {
		panic(err)
	}

	assetDir := path.Join(themeDir, themeConfig.Asset.Dir)
	themeAssets, err := filepath.Glob(assetDir + "/*/**")

	if err != nil {
		panic(err)
	}

	for _, themeAsset := range themeAssets {
		themeAssetName := strings.Trim(strings.Replace(themeAsset, assetDir, "", -1), "/")
		themeAssetType := path.Ext(themeAssetName)

		if themeAssetType == ".css" {
			themeAssetType = "css"
		} else if themeAssetType == ".js" {
			themeAssetType = "js"
		} else {
			themeAssetType = "other"
		}

		AppendAsset(&StaticAsset{
			Name:          themeAssetName,
			Path:          themeAsset,
			Type:          themeAssetType,
			DisableInline: themeAssetType == "other" || utils.SliceContains(themeConfig.Asset.DisableInline, themeAssetName),
			DisableMinify: themeAssetType == "other" || utils.SliceContains(themeConfig.Asset.DisableMinify, themeAssetName),
		})
	}
}

func AppendAsset(asset *StaticAsset) {
	for _, a := range assets {
		if a.Name == asset.Name {
			return
		}
	}

	assets = append(assets, asset)
}

func All() []*StaticAsset {
	return assets
}

func decodeBase64(input string) string {
	s, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(s)
}

func getStaticFile(assetName, assetType string) *StaticAsset {
	for _, file := range assets {
		if file.Name == assetName && file.Type == assetType {
			return file
		}
	}

	return &StaticAsset{}
}

func CssFile(assetName string) string {
	assetFile := getStaticFile(assetName, "css")

	if !config.DEVELOPMENT && !assetFile.DisableInline && assetFile.Content != "" {
		return fmt.Sprintf(`<style type="text/css" data-file="%s">%s</style>`, assetName, assetFile.Content)
	}

	return fmt.Sprintf(
		`<link rel="stylesheet" href="%s" />`,
		utils.Url(path.Join("/assets", assetName)),
	)
}

func JsFile(assetName string) string {
	assetFile := getStaticFile(assetName, "js")

	if !config.DEVELOPMENT && !assetFile.DisableInline && assetFile.Content != "" {
		return fmt.Sprintf(`<script charset="utf-8" data-file="%s">%s</script>`, assetName, assetFile.Content)
	}

	return fmt.Sprintf(
		`<script charset="utf-8" src="%s"></script>`,
		utils.Url(path.Join("/assets", assetName)),
	)
}

func OtherFile(assetName string) string {
	return utils.Url(path.Join("/assets", assetName))
}
