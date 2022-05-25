package asset

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/test"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/stretchr/testify/assert"
)

func checkBundledAssets() {
	if _, err := os.Stat("./bundled.js.go"); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	if _, err := os.Stat("./bundled.css.go"); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	if _, err := os.Stat("./bundled.other.go"); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
}

func TestLoadWithBundledAssets(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	Load(themeDir, true)

	assert.Equal(t, true, themeAssetsBundled)
	assert.Equal(t, 14, len(bundledAssets))
	assert.Equal(t, len(assets), len(bundledAssets))
}

func TestLoadWithoutBundledAssets(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	config.DEVELOPMENT = false
	themeAssetsBundled = false
	assets = []*StaticAsset{}
	Load(themeDir, true)
	assert.Equal(t, true, true)
}

func TestCssFileLink(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	config.DEVELOPMENT = true
	Load(themeDir, true)
	assert.Equal(t, "<link rel=\"stylesheet\" href=\"/assets/css/style.css\" />", CssFile("css/style.css"))
}

func TestCssFileInline(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	config.DEVELOPMENT = false
	assetName := "css/inline-file.css"

	bundledAssets = append(bundledAssets, &StaticAsset{Name: assetName, Content: "body{color:red}", Type: "css"})
	Load(themeDir, true)
	styleFile := getStaticFile(assetName, "css")
	assert.Equal(t, fmt.Sprintf(`<style type="text/css" data-file="%s">%s</style>`, assetName, styleFile.Content), CssFile(assetName))

	bundledAssets = utils.SliceFilter(bundledAssets, func(asset *StaticAsset) bool {
		return asset.Name != assetName
	})
}

func TestJsFileLink(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	config.DEVELOPMENT = true
	Load(themeDir, true)
	assert.Equal(t, "<script charset=\"utf-8\" src=\"/assets/js/main.js\"></script>", JsFile("js/main.js"))
}

func TestJsFileInline(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	config.DEVELOPMENT = false
	assetName := "js/main.js"
	Load(themeDir, true)
	jsFile := getStaticFile(assetName, "js")
	assert.Equal(t, fmt.Sprintf(`<script charset="utf-8" data-file="%s">%s</script>`, assetName, jsFile.Content), JsFile(assetName))
}

func TestOtherFileLink(t *testing.T) {
	checkBundledAssets()
	themeDir := "../themes/default"
	config.DEVELOPMENT = true
	Load(themeDir, true)
	assert.Equal(t, utils.Url("/assets/images/sample-image-file.png"), OtherFile("images/sample-image-file.png"))
}

func TestPanicInvalidThemeDir(t *testing.T) {
	defer test.RecoverPanic(t, "open ../themes/notfound/theme.json: no such file or directory", "theme error")
	themeDir := "../themes/notfound"
	Load(themeDir, true)
}

func TestAssetFileNotFound(t *testing.T) {
	themeDir := "../themes/default"
	Load(themeDir, true)
	file := getStaticFile("css/notfound.css", "css")
	assert.Equal(t, "", file.Name)
}

func TestAssetFunctions(t *testing.T) {
	themeDir := "../themes/default"
	Load(themeDir, true)
	AppendAsset(&StaticAsset{Name: "css/style.css", Content: "", Type: "css"})
	assert.Equal(t, 14, len(assets))

	AppendAsset(&StaticAsset{Name: "css/style1.css", Content: "", Type: "css"})
	assert.Equal(t, 15, len(All()))
	assert.Equal(t, "", decodeBase64("invalid base64 string"))

	bundledAssets = append(bundledAssets, &StaticAsset{Name: "css/append_file.css", Content: "", Type: "css"})
	Load(themeDir, true)
	assert.Equal(t, len(bundledAssets), len(assets))

	themeAssetsBundled = true
	config.DEVELOPMENT = false
	assets = []*StaticAsset{}
	bundledAssets = []*StaticAsset{}
	Load(themeDir, false)
	assert.Equal(t, 0, len(assets))
}

func TestPanicInvalidThemeFile(t *testing.T) {
	defer test.RecoverPanic(t, "unexpected end of JSON input", "theme error")
	themeDir, err := ioutil.TempDir("../../private/tmp", "theme-")

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(themeDir)

	themeConfig := ``

	if err := os.WriteFile(themeDir+"/theme.json", []byte(themeConfig), 0644); err != nil {
		panic(err)
	}

	Load(themeDir, true)
}

func TestAssets(t *testing.T) {
	defer test.RecoverPanic(t, "syntax error in pattern", "theme error")
	themeDir, err := ioutil.TempDir("../../private/tmp", "theme-")

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(themeDir)

	themeConfig := `{
		"asset": {
			"dir": "[-]",
			"disable_inline": [
				"js/highlight-11.5.0.min.js"
			],
			"disable_minify": [
				"js/highlight-11.5.0.min.js"
			]
		}
	}
	`

	if err := os.WriteFile(themeDir+"/theme.json", []byte(themeConfig), 0644); err != nil {
		panic(err)
	}

	Load(themeDir, true)
}
