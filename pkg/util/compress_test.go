package util

import (
	"os"
	"testing"
)

const (
	testZipFileName = "test.zip"
)

func TestArchive(t *testing.T) {
	type args struct {
		dirPath     string
		zipFileName string
	}
	os.MkdirAll("./test", 0755)
	defer os.RemoveAll("./test")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Archive successfully",
			args: args{
				dirPath:     "./test",
				zipFileName: testZipFileName,
			},
			wantErr: false,
		},
		{
			name: "Archive error, since there is no file/dir",
			args: args{
				dirPath:     "./test1",
				zipFileName: testZipFileName,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Archive(tt.args.dirPath, tt.args.zipFileName); (err != nil) != tt.wantErr {
				t.Errorf("Archive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnarchive(t *testing.T) {
	outputPath := "./test"
	defer os.RemoveAll(outputPath)

	type args struct {
		zipFileName string
		outputPath  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "UnArchive successfully",
			args: args{
				zipFileName: testZipFileName,
				outputPath:  outputPath,
			},
			wantErr: false,
		},
		{
			name: "UnArchive error, since there is no zip file",
			args: args{
				zipFileName: "./nozip.zip",
				outputPath:  outputPath,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unarchive(tt.args.zipFileName, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("Unarchive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
