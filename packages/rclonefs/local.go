package rclonefs

import (
	"context"
	"os"

	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/rclone/rclone/backend/local"
	"github.com/rclone/rclone/fs/config/configmap"
)

type RcloneLocal struct {
	*BaseRcloneDisk
	Root      string        `json:"root"`
	BaseUrl   string        `json:"base_url"`
	BaseUrlFn func() string `json:"-"`
}

type RcloneLocalConfig struct {
	Name      string        `json:"name"`
	Root      string        `json:"root"`
	BaseUrl   string        `json:"base_url"`
	BaseUrlFn func() string `json:"-"`
}

func NewLocal(cfg *RcloneLocalConfig) fs.FSDisk {
	rl := &RcloneLocal{
		BaseRcloneDisk: &BaseRcloneDisk{
			DiskName: cfg.Name,
		},
		Root:      cfg.Root,
		BaseUrl:   cfg.BaseUrl,
		BaseUrlFn: cfg.BaseUrlFn,
	}

	if err := os.MkdirAll(cfg.Root, os.ModePerm); err != nil {
		panic(err)
	}

	cfgMap := configmap.New()
	cfgMap.Set("root", rl.Root)
	fsDriver, err := local.NewFs(context.Background(), rl.DiskName, rl.Root, cfgMap)

	if err != nil {
		panic(err)
	}

	rl.Fs = fsDriver

	return rl
}

func (r *RcloneLocal) Url(filepath string) string {
	if r.BaseUrlFn != nil {
		return r.BaseUrlFn() + "/" + filepath
	}
	return r.BaseUrl + "/" + filepath
}

func (r *RcloneLocal) Delete(ctx context.Context, filepath string) error {
	obj, err := r.Fs.NewObject(ctx, filepath)

	if err != nil {
		return err
	}

	return obj.Remove(ctx)
}
