package rclonefs

import (
	"github.com/ngocphuongnb/tetua/app/fs"
)

func NewFromConfig(config *fs.StorageConfig) []fs.FSDisk {
	var disks []fs.FSDisk

	for _, diskConfig := range config.DiskConfigs {
		switch diskConfig.Driver {
		case "s3":
			disks = append(disks, NewS3(&RcloneS3Config{
				Name:            diskConfig.Name,
				Root:            diskConfig.Root,
				Provider:        diskConfig.Provider,
				Bucket:          diskConfig.Bucket,
				Region:          diskConfig.Region,
				Endpoint:        diskConfig.Endpoint,
				AccessKeyID:     diskConfig.AccessKeyID,
				SecretAccessKey: diskConfig.SecretAccessKey,
				BaseUrl:         diskConfig.BaseUrl,
				ACL:             diskConfig.ACL,
			}))
		case "local":
			disks = append(disks, NewLocal(&RcloneLocalConfig{
				Name:      diskConfig.Name,
				Root:      diskConfig.Root,
				BaseUrl:   diskConfig.BaseUrl,
				BaseUrlFn: diskConfig.BaseUrlFn,
			}))
		}
	}

	return disks
}
