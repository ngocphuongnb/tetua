package fs_test

import (
	"testing"

	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/stretchr/testify/assert"
)

type Disk2 struct {
	*mock.Disk
}

func (d *Disk2) Name() string {
	return "disk_mock2"
}

func TestNewNoDisk(t *testing.T) {
	defer utils.RecoverTestPanic(t, "No disk found")
	fs.New("", nil)
	assert.Equal(t, 1, 1)
}

func TestNewNoDefaultDisk(t *testing.T) {
	defer utils.RecoverTestPanic(t, "No default disk found")
	fs.New("", []fs.FSDisk{&mock.Disk{}})
	assert.Equal(t, 1, 1)
}

func TestNew(t *testing.T) {
	disk1 := &mock.Disk{}
	disk2 := &Disk2{}
	fs.New("disk_mock", []fs.FSDisk{disk1, disk2})
	assert.Equal(t, disk1, fs.Disk())
	assert.Equal(t, disk2, fs.Disk("disk_mock2"))
	assert.Equal(t, nil, fs.Disk("disk_mock_test"))

}
