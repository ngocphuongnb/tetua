package utils_test

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestSanitizePlainText(t *testing.T) {
	assert.Equal(t, "test", utils.SanitizePlainText("<script>alert('test')</script><p>test</p>"))
}

func TestSanitizeMarkdown(t *testing.T) {
	assert.Equal(t,
		"# Title\n**Code**\n\n```<?php echo 1; ?>```\n<iframe src=\"https://www.youtube.com/embed/dQw4w9WgXcQ\" frameborder=\"0\" allowfullscreen=\"\" sandbox=\"\"></iframe>\n<iframe frameborder=\"0\" allowfullscreen=\"\" sandbox=\"\"></iframe>",
		utils.SanitizeMarkdown("# Title\n**Code**\n\n```<?php echo 1; ?>```\n<iframe src=\"https://www.youtube.com/embed/dQw4w9WgXcQ\" frameborder=\"0\" allowfullscreen=\"\" sandbox=\"\"></iframe>\n<iframe src=\"https://danger.local\" frameborder=\"0\" allowfullscreen=\"\" sandbox=\"\"></iframe>"),
	)
}

func TestExtractContent(t *testing.T) {
	title, content := utils.ExtractContent("# Title\n## Content")
	assert.Equal(t, "Title", title)
	assert.Equal(t, "## Content", content)
}

func TestMarkdownToHtml(t *testing.T) {
	html, err := utils.MarkdownToHtml("# Title\n## Content")
	assert.NoError(t, err)
	assert.Equal(t, "<h1 id=\"title\">Title</h1>\n<h2 id=\"content\">Content</h2>\n", html)
}

func TestGenerateHash(t *testing.T) {
	hash, err := utils.GenerateHash("")
	assert.Equal(t, "hash: input cannot be empty", err.Error())
	assert.Equal(t, "", hash)

	hash, err = utils.GenerateHash("input")
	assert.NoError(t, err)
	assert.Equal(t, true, hash != "")

	err = utils.CheckHash("input", "")
	assert.Equal(t, "invalid hash", err.Error())

	err = utils.CheckHash("input", "$argon2id$v=a$m=a,t=a,p=a$OamJg0HQDzROGA8uxX6QtA$B3Ei9eglsBPiyCrxLGdqV2KeYIRTpfsvBLH4GSOEw8M")
	assert.Equal(t, "expected integer", err.Error())

	err = utils.CheckHash("input", "$argon2id$v=19$m=65536,t=3,p=4$invalid_salt$B3Ei9eglsBPiyCrxLGdqV2KeYIRTpfsvBLH4GSOEw8M")
	assert.Equal(t, "illegal base64 data at input byte 7", err.Error())

	err = utils.CheckHash("input", "$argon2id$v=19$m=65536,t=3,p=4$OamJg0HQDzROGA8uxX6QtA$invalid_hash")
	assert.Equal(t, "illegal base64 data at input byte 7", err.Error())

	err = utils.CheckHash("input2", "$argon2id$v=19$m=65536,t=3,p=4$OamJg0HQDzROGA8uxX6QtA$B3Ei9eglsBPiyCrxLGdqV2KeYIRTpfsvBLH4GSOEw8M")
	assert.Equal(t, "invalid hash", err.Error())

	err = utils.CheckHash("input", hash)
	assert.NoError(t, err)
}

func TestGetFunctionName(t *testing.T) {
	assert.Equal(t, true, strings.HasSuffix(utils.GetFunctionName(TestGenerateHash), "utils_test.TestGenerateHash"))
	defer utils.RecoverTestPanic(t, "reflect: call of reflect.Value.Pointer on zero Value", "nil input")
	assert.Equal(t, "", utils.GetFunctionName(nil))
}

func TestSliceOverlap(t *testing.T) {
	assert.Equal(t, []int{2, 3}, utils.SliceOverlap([]int{1, 2, 3}, []int{2, 3, 4}))
	assert.Equal(t, []int{}, utils.SliceOverlap([]int{1, 2, 3}, []int{4, 5, 6}))
}

func TestSliceFilter(t *testing.T) {
	assert.Equal(t, []int{2, 3}, utils.SliceFilter([]int{1, 2, 3}, func(i int) bool { return i > 1 }))
}

func TestSliceMap(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, utils.SliceMap([]int{1, 2, 3}, func(i int) string { return strconv.Itoa(i) }))
}

func TestRepeat(t *testing.T) {
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, utils.Repeat(1, 10))
}

func TestSliceAppendIfNotExists(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, utils.SliceAppendIfNotExists([]int{1, 2, 3}, 4, func(i int) bool {
		return i == 4
	}))
	assert.Equal(t, []int{1, 2, 3}, utils.SliceAppendIfNotExists([]int{1, 2, 3}, 2, func(i int) bool {
		return i == 2
	}))
}

func TestRecoverPanic(t *testing.T) {
	defer utils.RecoverTestPanic(t, "panic", "panic")
	panic("panic")
}

func TestGetRootDir(t *testing.T) {
	assert.Equal(t, true, strings.HasSuffix(utils.GetRootDir(), "go/src/testing"))
}

func TestCreateTestDir(t *testing.T) {
	dir := utils.CreateTestDir("testing-dir-")
	defer os.RemoveAll(dir)
	assert.Equal(t, true, strings.Contains(dir, "testing-dir-"))
}

func TestGetStructField(t *testing.T) {
	type testStruct struct {
		Field1 string
		Field2 int
	}
	var s testStruct = testStruct{
		Field1: "field1 value",
		Field2: 2,
	}
	assert.Equal(t, "field1 value", utils.GetStructField(s, "Field1").String())
	assert.Equal(t, 2, int(utils.GetStructField(s, "Field2").Int()))
}

func TestFirstError(t *testing.T) {
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	assert.Equal(t, nil, utils.FirstError(nil))
	assert.Equal(t, err1, utils.FirstError(nil, err1))
	assert.Equal(t, err2, utils.FirstError(nil, err2, err1))
}
