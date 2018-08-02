package dropbox

import (
	"io"
	"os"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

var (
	// Dropbox alow maximum 150Mb (150 * 1048576) per session
	// In our tests it performed better with 1mb.
	chunkSize int64 = 1048576
)

// Node contains metadata to files and folders
type Node struct {
	IsFolder       bool
	Name           string
	Size           uint64
	Rev            string
	ServerModified time.Time
}

// NewConfig return the dropbox config
func NewConfig(token string) (config dropbox.Config) {
	config = dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogOff, // if needed, set the desired logging level. Default is off
	}
	return
}

// List files and folders in the path
// do not use "/" to root folder instead use ""
func List(config dropbox.Config, path string) (nodes []Node, err error) {
	f := files.New(config)
	lfa := files.NewListFolderArg(path)
	lfr, err := f.ListFolder(lfa)
	if err != nil {
		return
	}
	for _, v := range lfr.Entries {
		var n Node
		switch fm := v.(type) {
		case *files.FileMetadata:
			n = parseFileMetadata(fm)
		case *files.FolderMetadata:
			n = parseFolderMetadata(fm)
		}
		nodes = append(nodes, n)
	}
	return
}

// parseFolderMetadata return a folder node from FolderMetadata
func parseFolderMetadata(e *files.FolderMetadata) (n Node) {
	n.IsFolder = true
	n.Name = e.Name
	return
}

// parseFileMetadata return a file node from FileMetadata
func parseFileMetadata(e *files.FileMetadata) (n Node) {
	n.Name = e.Name
	n.Rev = e.Rev
	n.ServerModified = e.ServerModified
	n.Size = e.Size
	return
}

// Upload file to Dropbox
func Upload(config dropbox.Config, fromFile, toFile string) (err error) {
	info, err := os.Lstat(fromFile)
	if err != nil {
		return err
	}

	f, err := os.Open(fromFile)
	if err != nil {
		return err
	}
	defer f.Close() // nolint
	ci := files.NewCommitInfo(toFile)
	ci.ClientModified = info.ModTime().UTC().Round(time.Second)
	client := files.New(config)
	if info.Size() < chunkSize {
		_, err = client.Upload(ci, f)
		return err
	}

	// create a reader with the first portion of the file
	r := io.LimitReader(f, chunkSize)

	// start the upload session
	s, err := client.UploadSessionStart(files.NewUploadSessionStartArg(), r)
	if err != nil {
		return err
	}
	var uploaded = chunkSize // bytes sended

	// loop sending file
	for (info.Size() - uploaded) > chunkSize {
		cursor := files.NewUploadSessionCursor(s.SessionId, uint64(uploaded))
		arg := files.NewUploadSessionAppendArg(cursor)
		r = io.LimitReader(f, chunkSize)
		err = client.UploadSessionAppendV2(arg, r)
		if err != nil {
			return err
		}
		uploaded += chunkSize
	}

	// send last file chunk
	cursor := files.NewUploadSessionCursor(s.SessionId, uint64(uploaded))
	arg := files.NewUploadSessionFinishArg(cursor, ci)

	// change commit mode to overwrite (defaulti is add)
	arg.Commit.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}}

	_, err = client.UploadSessionFinish(arg, f)
	return err
}

// Download files from dropbox
func Download(config dropbox.Config, fromFile, toFile string) (err error) {
	client := files.New(config)
	da := files.NewDownloadArg(fromFile)
	_, src, err := client.Download(da)
	if err != nil {
		return
	}
	defer src.Close() // nolint
	dst, err := os.Create(toFile)
	if err != nil {
		return
	}
	defer dst.Close() // nolint
	_, err = io.Copy(dst, src)
	return
}
