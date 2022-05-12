package mock

import (
	"context"
	"errors"
	"io"
	"mime/multipart"

	"github.com/ngocphuongnb/tetua/app/fs"
)

type Disk struct {
}

func (d *Disk) Name() string {
	return "disk_mock"
}

func (d *Disk) Url(path string) string {
	return "/disk_mock/" + path
}

func (d *Disk) Delete(ctx context.Context, path string) error {
	if path == "/delete/error" {
		return errors.New("Delete file error")
	}
	return nil
}

func (d *Disk) Put(ctx context.Context, in io.Reader, size int64, mime, dst string) (*fs.FileInfo, error) {
	return nil, nil
}

func (d *Disk) PutMultipart(ctx context.Context, m *multipart.FileHeader, dsts ...string) (*fs.FileInfo, error) {
	if m.Filename == "error.jpg" {
		return nil, errors.New("PutMultipart error")
	}

	mime := ""

	if mimes := m.Header["Content-Type"]; len(mimes) > 0 {
		mime = mimes[0]
	}

	return &fs.FileInfo{
		Disk: d.Name(),
		Path: m.Filename,
		Type: mime,
		Size: 100,
	}, nil
}
