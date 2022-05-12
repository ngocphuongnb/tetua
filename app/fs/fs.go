package fs

import (
	"context"
	"io"
	"mime/multipart"
)

var defaultStorageDisk FSDisk = nil

type FSDisk interface {
	Name() string
	Url(filepath string) string
	Delete(ctx context.Context, filepath string) error
	Put(ctx context.Context, in io.Reader, size int64, mime, dst string) (*FileInfo, error)
	PutMultipart(ctx context.Context, m *multipart.FileHeader, dsts ...string) (*FileInfo, error)
}

type FileInfo struct {
	Disk string `json:"disk,omitempty"`
	Path string `json:"path,omitempty"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
}

type DiskConfig struct {
	Name            string        `json:"name"`
	Driver          string        `json:"driver"`
	Root            string        `json:"root"`
	BaseUrl         string        `json:"base_url"`
	BaseUrlFn       func() string `json:"-"`
	Provider        string        `json:"provider"`
	Endpoint        string        `json:"endpoint"`
	Region          string        `json:"region"`
	Bucket          string        `json:"bucket"`
	AccessKeyID     string        `json:"access_key_id"`
	SecretAccessKey string        `json:"secret_access_key"`
	ACL             string        `json:"acl"`
}

type StorageConfig struct {
	DefaultDisk string        `json:"default_disk"`
	DiskConfigs []*DiskConfig `json:"disks"`
}

var fsDisks []FSDisk

func New(defaultDisk string, disks []FSDisk) {
	fsDisks = append(fsDisks, disks...)

	if len(fsDisks) == 0 {
		panic("No disk found")
	}

	for _, disk := range fsDisks {
		if disk.Name() == defaultDisk {
			defaultStorageDisk = disk
			break
		}
	}

	if defaultStorageDisk == nil {
		panic("No default disk found")
	}
}

func Disk(names ...string) FSDisk {
	if len(names) == 0 {
		return defaultStorageDisk
	}

	for _, disk := range fsDisks {
		if disk.Name() == names[0] {
			return disk
		}
	}

	return nil
}
