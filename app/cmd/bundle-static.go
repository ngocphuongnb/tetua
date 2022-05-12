package cmd

import (
	"encoding/base64"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"

	"github.com/ngocphuongnb/tetua/app/asset"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

var m = minify.New()

var (
	cssBundledContent   = ""
	jsBundledContent    = ""
	otherBundledContent = ""

	cssBundledFile   = ""
	jsBundledFile    = ""
	otherBundledFile = ""
)

func reset() {
	cssBundledContent = "package asset\n"
	jsBundledContent = "package asset\n"
	otherBundledContent = "package asset\n"
	cssBundledFile = path.Join(config.WD, "app/asset/bundled.css.go")
	jsBundledFile = path.Join(config.WD, "app/asset/bundled.js.go")
	otherBundledFile = path.Join(config.WD, "app/asset/bundled.other.go")
}

func init() {
	reset()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
}

func bundleAssets() (err error) {
	cssBundledContent += "func init() {\n"
	cssBundledContent += "themeAssetsBundled = true\n"

	jsBundledContent += "func init() {\n"
	jsBundledContent += "themeAssetsBundled = true\n"

	otherBundledContent += "func init() {\n"
	otherBundledContent += "themeAssetsBundled = true\n"

	appAssets := asset.All()

	for _, assetFile := range appAssets {
		assetContent := ""
		assetFilePath := assetFile.Path

		b, err := ioutil.ReadFile(assetFilePath)
		if err != nil {
			return err
		}

		assetContent += string(b)

		if !assetFile.DisableMinify {
			if assetFile.Type == "css" {
				assetContent, err = m.String("text/css", assetContent)
			}

			if assetFile.Type == "js" {
				assetContent, err = m.String("application/javascript", assetContent)
			}

			if err != nil {
				return err
			}
		}

		assetContent = base64.StdEncoding.EncodeToString([]byte(assetContent))
		assetBundledContent := "bundledAssets = append(bundledAssets, &StaticAsset{\n"
		assetBundledContent += "Type:\"" + assetFile.Type + "\",\n"
		assetBundledContent += "Name:\"" + assetFile.Name + "\",\n"
		assetBundledContent += "Path:\"" + assetFile.Path + "\",\n"
		assetBundledContent += fmt.Sprintf("DisableMinify:%t,\n", assetFile.DisableMinify)
		assetBundledContent += fmt.Sprintf("DisableInline:%t,\n", assetFile.DisableInline)
		assetBundledContent += fmt.Sprintf("Content: decodeBase64(\"%s\"),\n", assetContent)
		assetBundledContent += "})\n"

		if assetFile.Type == "css" {
			cssBundledContent += assetBundledContent
		} else if assetFile.Type == "js" {
			jsBundledContent += assetBundledContent
		} else {
			otherBundledContent += assetBundledContent
		}

	}
	cssBundledContent += "}"
	jsBundledContent += "}"
	otherBundledContent += "}"

	cssBundledContentBytes, err := format.Source([]byte(cssBundledContent))
	if err != nil {
		return err
	}

	jsBundledContentBytes, err := format.Source([]byte(jsBundledContent))
	if err != nil {
		return err
	}

	otherBundledContentBytes, err := format.Source([]byte(otherBundledContent))
	if err != nil {
		return err
	}

	if err = os.WriteFile(cssBundledFile, cssBundledContentBytes, 0644); err != nil {
		return err
	}

	if err = os.WriteFile(jsBundledFile, jsBundledContentBytes, 0644); err != nil {
		return err
	}

	if err = os.WriteFile(otherBundledFile, otherBundledContentBytes, 0644); err != nil {
		return err
	}

	return nil
}

func BundleStaticAssets() (err error) {
	reset()
	return bundleAssets()
}
