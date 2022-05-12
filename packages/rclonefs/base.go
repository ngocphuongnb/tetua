package rclonefs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/utils"

	rclonefs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/object"
)

var filenameRemoveCharsRegexp = regexp.MustCompile(`[^a-zA-Z0-9_\-\.]`)
var dashRegexp = regexp.MustCompile(`\-+`)
var allowedMimes = []string{
	"text/xml",
	"text/xml; charset=utf-8",
	"image/svg+xml",
	"image/jpeg",
	"image/pjpeg",
	"image/png",
	"image/gif",
	"image/x-icon",
	"application/pdf",
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/powerpoint",
	"application/x-mspowerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/mspowerpoint",
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.slideshow",
	"application/vnd.oasis.opendocument.text",
	"application/excel",
	"application/vnd.ms-excel",
	"application/x-excel",
	"application/x-msexcel",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	// "application/octet-stream",
	"audio/mpeg3",
	"audio/x-mpeg-3",
	"video/x-mpeg",
	"audio/m4a",
	"audio/ogg",
	"audio/wav",
	"audio/x-wav",
	"video/mp4",
	"video/x-m4v",
	"video/quicktime",
	"video/x-ms-asf",
	"video/x-ms-wmv",
	"application/x-troff-msvideo",
	"video/avi",
	"video/msvideo",
	"video/x-msvideo",
	"audio/mpeg",
	"video/mpeg",
	"video/ogg",
	"video/3gpp",
	"audio/3gpp",
	"video/3gpp2",
	"audio/3gpp2",
}

type BaseRcloneDisk struct {
	rclonefs.Fs
	DiskName string `json:"name"`
	Root     string
}

func (r *BaseRcloneDisk) Name() string {
	return r.DiskName
}

func (r *BaseRcloneDisk) Put(ctx context.Context, reader io.Reader, size int64, mime, dst string) (*fs.FileInfo, error) {
	objectInfo := object.NewStaticObjectInfo(
		dst,
		time.Now(),
		size,
		true,
		nil,
		nil,
	)

	rs, err := r.Fs.Put(ctx, reader, objectInfo)

	if err != nil {
		return nil, err
	}

	return &fs.FileInfo{
		Disk: r.DiskName,
		Path: dst,
		Type: mime,
		Size: int(rs.Size()),
	}, nil
}

func (r *BaseRcloneDisk) PutMultipart(ctx context.Context, m *multipart.FileHeader, dsts ...string) (*fs.FileInfo, error) {
	f, err := m.Open()

	if err != nil {
		return nil, err
	}

	fileHeader := make([]byte, 512)

	if _, err := f.Read(fileHeader); err != nil {
		return nil, err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	dst := ""
	mime := http.DetectContentType(fileHeader)

	if !utils.SliceContains(allowedMimes, strings.ToLower(mime)) {
		return nil, errors.New("file type is not allowed")
	}

	if len(dsts) > 0 {
		dst = dsts[0]
	} else {
		dst = r.UploadFilePath(m.Filename)
	}

	return r.Put(ctx, f, m.Size, mime, dst)
}

func (r *BaseRcloneDisk) UploadFilePath(filename string) string {
	now := time.Now()
	filename = filenameRemoveCharsRegexp.ReplaceAllString(filename, "-")
	filename = dashRegexp.ReplaceAllString(filename, "-")
	filename = strings.ReplaceAll(filename, "-.", ".")
	return path.Join(
		r.Root,
		strconv.Itoa(now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		fmt.Sprintf("%d_%s", now.UnixMicro(), filename),
	)
}
