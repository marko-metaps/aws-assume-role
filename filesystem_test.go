package main

import (
	"os"
	"testing"
	"time"
)

type MockFileInfo struct {
	perm os.FileMode
}

func (m MockFileInfo) Name() string       { return "" }
func (m MockFileInfo) Size() int64        { return 0 }
func (m MockFileInfo) Mode() os.FileMode  { return m.perm }
func (m MockFileInfo) ModTime() time.Time { return time.Time{} }
func (m MockFileInfo) IsDir() bool        { return false }
func (m MockFileInfo) Sys() interface{}   { return nil }

type MockFileSystem struct {
	HomeDirErr  error
	HomeDirPath string
	StatErr     error
	FileInfo    MockFileInfo
}

func (m MockFileSystem) UserHomeDir() (string, error) {
	return m.HomeDirPath, m.HomeDirErr
}

func (m MockFileSystem) Stat(path string) (os.FileInfo, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	return m.FileInfo, nil
}

func TestCheckCredentialFile(t *testing.T) {
	tests := []struct {
		name        string
		fs          FileSystem
		expectError bool
	}{
		{
			name: "file does not exist",
			fs: MockFileSystem{
				HomeDirPath: "/dummy/home",
				StatErr:     os.ErrNotExist,
			},
			expectError: true,
		},
		{
			name: "no write permission",
			fs: MockFileSystem{
				HomeDirPath: "/dummy/home",
				FileInfo:    MockFileInfo{perm: 0444},
			},
			expectError: true,
		},
		{
			name: "has write permission",
			fs: MockFileSystem{
				HomeDirPath: "/dummy/home",
				FileInfo:    MockFileInfo{perm: 0644},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCredentialFile(tt.fs)
			if (err != nil) != tt.expectError {
				t.Errorf("checkCredentialFile() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
