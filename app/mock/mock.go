package mock

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"strings"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/server"
	fiber "github.com/ngocphuongnb/tetua/packages/fiberserver"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func CreateLogger(silences ...bool) *MockLogger {
	silences = append(silences, false)
	mockLogger := &MockLogger{
		Silence: silences[0],
	}
	logger.New(mockLogger)
	return mockLogger
}

func CreateServer() server.Server {
	if logger.Get() == nil {
		CreateLogger()
	}
	s := fiber.New(fiber.Config{JwtSigningKey: "sesj5JYrRxrB2yUWkBFM7KKWCY2ykxBw"})
	s.Use(auth.AssignUserInfo)
	s.Use(auth.Check)

	return s
}

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func CreateUploadRequest(method, uri, fieldName, fileName string) *http.Request {
	uploadFilePath := fileName
	requestBody := new(bytes.Buffer)
	writer := multipart.NewWriter(requestBody)
	contentType := "application/octet-stream"

	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		contentType = "image/jpeg"
	}

	if strings.HasSuffix(fileName, ".png") {
		contentType = "image/png"
	}

	header := make(textproto.MIMEHeader)
	header.Set(
		"Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldName), escapeQuotes(fileName)),
	)
	header.Set("Content-Type", contentType)
	part, _ := writer.CreatePart(header)

	sampleFile, _ := os.Open(uploadFilePath)
	io.Copy(part, sampleFile)
	writer.Close()

	req := httptest.NewRequest(method, uri, requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func Request(s server.Server, method, uri string, headers ...map[string]string) (string, *http.Response) {
	req := httptest.NewRequest(method, uri, nil)
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}
	resp, _ := s.Test(req)
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), resp
}

func SendRequest(s server.Server, req *http.Request) (string, *http.Response) {
	resp, _ := s.Test(req)
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), resp
}

func GetRequest(s server.Server, uri string, headers ...map[string]string) (string, *http.Response) {
	return Request(s, "GET", uri, headers...)
}

func PostRequest(s server.Server, uri string, headers ...map[string]string) (string, *http.Response) {
	return Request(s, "POST", uri, headers...)
}
