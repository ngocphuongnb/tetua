package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	h "html"

	"github.com/microcosm-cc/bluemonday"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Iterations uint32
	Memory     uint32
	KeyLen     uint32
	Threads    uint8
}

// var ugcPolicy = bluemonday.UGCPolicy()
var stripTagsPolicy = bluemonday.StripTagsPolicy()
var markdownPolicy = bluemonday.StripTagsPolicy()
var md goldmark.Markdown

var iframeAllowHosts = []string{
	"www.youtube.com",
	"codesandbox.io",
	"gist.github.com",
	"instagram.com",
	"twitter.com",
	"twitch.tv",
	"vimeo.com",
	"codepen.io",
	"glitch.com",
	"jsbin.com",
	"jsfiddle.net",
	"repl.it",
	"reddit.com",
	"slideshare.net",
	"soundcloud.com",
	"stackblitz.com",
}

func init() {
	markdownPolicy.AllowIFrames()
	markdownPolicy.AllowElements("div")
	markdownPolicy.AllowElements("php")
	markdownPolicy.AllowElements("?php")
	markdownPolicy.AllowElements("iframe")
	markdownPolicy.
		AllowURLSchemeWithCustomPolicy("https", func(url *url.URL) (allowUrl bool) {
			return SliceContains(iframeAllowHosts, url.Host)
		}).
		AllowAttrs("src", "frameborder", "allowfullscreen", "width", "height", "allow").
		OnElements("iframe")

	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
}

func SanitizePlainText(html string) string {
	return stripTagsPolicy.Sanitize(html)
}

func SanitizeMarkdown(html string) string {
	html = strings.ReplaceAll(html, "<?php", "__php_open_tag__")
	html = h.UnescapeString(markdownPolicy.Sanitize(html))
	return strings.ReplaceAll(html, "__php_open_tag__", "<?php")
}

func ExtractContent(content string) (string, string) {
	content = SanitizeMarkdown(content)
	lines := strings.Split(content, "\n")
	name := strings.Trim(strings.Trim(lines[0], "#"), " ")
	content = strings.Join(lines[1:], "\n")

	return name, content
}

func MarkdownToHtml(content string) (string, error) {
	var buf bytes.Buffer

	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateHash(input string) (string, error) {
	if input == "" {
		return "", errors.New("hash: input cannot be empty")
	}

	salt := make([]byte, 16)
	cfg := &HashConfig{
		Iterations: 3,
		Memory:     64 * 1024,
		Threads:    4,
		KeyLen:     32,
	}

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(input), salt, cfg.Iterations, cfg.Memory, cfg.Threads, cfg.KeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, cfg.Memory, cfg.Iterations, cfg.Threads, b64Salt, b64Hash)

	return full, nil
}

func CheckHash(input, hash string) error {
	var err error
	var salt []byte
	var decodedHash []byte
	parts := strings.Split(hash, "$")
	cfg := &HashConfig{}

	if len(parts) != 6 {
		return errors.New("invalid hash")
	}

	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &cfg.Memory, &cfg.Iterations, &cfg.Threads); err != nil {
		return err
	}

	if salt, err = base64.RawStdEncoding.DecodeString(parts[4]); err != nil {
		return err
	}

	if decodedHash, err = base64.RawStdEncoding.DecodeString(parts[5]); err != nil {
		return err
	}

	cfg.KeyLen = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(input), salt, cfg.Iterations, cfg.Memory, cfg.Threads, cfg.KeyLen)
	valid := subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1

	if !valid {
		return errors.New("invalid hash")
	}

	return nil
}

func Encrypt(stringToEncrypt string, keys ...string) (string, error) {
	key := []byte(config.APP_KEY)

	if len(keys) > 0 {
		key = []byte(keys[0])
	}

	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func Decrypt(encryptedString string, keys ...string) (string, error) {
	key := []byte(config.APP_KEY)
	if len(keys) > 0 {
		key = []byte(keys[0])
	}

	enc, err := hex.DecodeString(encryptedString)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(enc) < aesGCM.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func Url(path string) string {
	appBase := strings.TrimRight(config.Setting("app_base_url"), "/")

	if path == "" {
		return appBase + "/"
	}

	path = strings.TrimLeft(path, "/")
	path = fmt.Sprintf("%s/%s", appBase, path)
	return path
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func SliceContains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func SliceOverlap[T comparable](slice1 []T, slice2 []T) []T {
	result := make([]T, 0)
	for _, e1 := range slice1 {
		if SliceContains(slice2, e1) {
			result = append(result, e1)
		}
	}
	return result
}

func SliceFilter[T comparable](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, e := range slice {
		if predicate(e) {
			result = append(result, e)
		}
	}
	return result
}

func SliceMap[T comparable, R comparable](slice []T, mapper func(T) R) []R {
	var result []R
	for _, e := range slice {
		result = append(result, mapper(e))
	}
	return result
}

func Repeat[T comparable](input T, time int) []T {
	var result []T
	var i = 0
	for i < time {
		result = append(result, input)
		i++
	}
	return result
}

func SliceAppendIfNotExists[T comparable](slice []T, newItem T, checkExists func(T) bool) []T {
	for _, s := range slice {
		if checkExists(s) {
			return slice
		}
	}
	slice = append(slice, newItem)
	return slice
}

// GetStructField returns the value of a struct field
func GetStructField(entity interface{}, field string) reflect.Value {
	r := reflect.ValueOf(entity)
	return reflect.Indirect(r).FieldByName(field)
}

// FirstError return the first error in a list of errors
func FirstError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
