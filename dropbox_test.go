package dropbox

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
)

func TestList(t *testing.T) {
	token := os.Getenv("DROPBOX_TOKEN")
	type args struct {
		config dropbox.Config
		path   string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		validate func([]Node) error
	}{
		{
			name:    "success",
			wantErr: false,
			args: args{
				config: NewConfig(token),
				path:   "",
			},
			validate: func(nodes []Node) (err error) {
				if len(nodes) == 0 {
					err = errors.New("more files are expected")
				}
				return
			},
		},
		{
			name:    "error folser path",
			wantErr: true,
			args: args{
				config: NewConfig(token),
				path:   "/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNodes, err := List(tt.args.config, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validate != nil {
				err = tt.validate(gotNodes)
				if err != nil {
					t.Error(err)
					return
				}
			}
		})
	}
}

func TestUpload(t *testing.T) {
	token := os.Getenv("DROPBOX_TOKEN")
	chunkSize = 1048576 / 4
	type args struct {
		config   dropbox.Config
		fromFile string
		toFile   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
			args: args{
				config:   NewConfig(token),
				fromFile: "fixtures/test.payload",
				toFile:   "/test.payload",
			},
		},
		{
			name:    "success small file",
			wantErr: false,
			args: args{
				config:   NewConfig(token),
				fromFile: "fixtures/test_small.payload",
				toFile:   "/test_small.payload",
			},
		},
		{
			name:    "error local file",
			wantErr: true,
			args: args{
				config:   NewConfig(token),
				fromFile: "",
				toFile:   "/test_small.payload",
			},
		},
		{
			name:    "error SessionStart",
			wantErr: true,
			args: args{
				config:   NewConfig(token),
				fromFile: "fixtures/test.payload",
				toFile:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Upload(tt.args.config, tt.args.fromFile, tt.args.toFile); (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func md5chksum(file, expected string) (err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	md5sum := fmt.Sprintf("%X", md5.Sum(b))
	if md5sum != expected {
		err = errors.New(fmt.Sprintf("expected %q but get %q", expected, md5sum))
	}
	return
}
func TestDownload(t *testing.T) {
	token := os.Getenv("DROPBOX_TOKEN")
	// Send files to test download
	err := Upload(NewConfig(token), "fixtures/test.payload", "/test.payload")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		config   dropbox.Config
		fromFile string
		toFile   string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		validate func(file string) (err error)
	}{
		{
			name:    "success",
			wantErr: false,
			args: args{
				config:   NewConfig(token),
				fromFile: "/test.payload",
				toFile:   "test.payload",
			},
			validate: func(file string) (err error) {
				return md5chksum(file, "B6D81B360A5672D80C27430F39153E2C")
			},
		},
		{
			name:    "error file not found",
			wantErr: true,
			args: args{
				config:   NewConfig(token),
				fromFile: "/not_found.payload",
				toFile:   "test.payload",
			},
			validate: func(file string) (err error) {
				return md5chksum(file, "B6D81B360A5672D80C27430F39153E2C")
			},
		},
		{
			name:    "error destination file error",
			wantErr: true,
			args: args{
				config:   NewConfig(token),
				fromFile: "/test.payload",
				toFile:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Download(tt.args.config, tt.args.fromFile, tt.args.toFile); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %q, wantErr %v", err, tt.wantErr)
			}
			if tt.validate != nil {
				err := tt.validate(tt.args.toFile)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	_ = os.RemoveAll("test.payload")
}
