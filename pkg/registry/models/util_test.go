package models

import (
	"net/http"
	"os"
	"path"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kleveross/ormb/pkg/ormb"
	ormbmock "github.com/kleveross/ormb/pkg/ormb/mock"
	. "github.com/onsi/ginkgo"
)

var (
	ormbClient ormb.Interface
)

func init() {
	ormbClient = ormbmock.NewMockInterface(gomock.NewController(GinkgoT()))

	ormbClient.(*ormbmock.MockInterface).EXPECT().Login(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(nil)
	ormbClient.(*ormbmock.MockInterface).EXPECT().Save(
		gomock.Any(),
		gomock.Any(),
	).Return(nil)
	ormbClient.(*ormbmock.MockInterface).EXPECT().Export(
		gomock.Any(),
		gomock.Any(),
	).Return(nil)
	ormbClient.(*ormbmock.MockInterface).EXPECT().Pull(
		gomock.Any(),
	).Return(nil)
	ormbClient.(*ormbmock.MockInterface).EXPECT().Push(
		gomock.Any(),
	).Return(nil)
	ormbClient.(*ormbmock.MockInterface).EXPECT().Remove(
		gomock.Any(),
	).Return(nil)
}

func Test_validateFileSize(t *testing.T) {
	requestRight := http.Request{
		Header: make(map[string][]string),
	}
	requestRight.Header["Content-Length"] = []string{"209715200"} // 209715200 = 100<<20+1

	requestError := http.Request{
		Header: make(map[string][]string),
	}
	requestError.Header["Content-Length"] = []string{"error"}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "validateFileSize error, file size is too large",
			args: args{
				r: &requestRight,
			},
			wantErr: true,
		},
		{
			name: "validateFileSize error, content length is not right",
			args: args{
				r: &requestError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFileSize(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("validateFileSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseContenRange(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int64
		want2   int64
		wantErr bool
	}{
		{
			name: "Parse content range successfully",
			args: args{
				s: "bytes 0-10485759/30319242",
			},
			want:    30319242,
			want1:   0,
			want2:   10485759,
			wantErr: false,
		},
		{
			name: "Parse content range error, since input is \"\"",
			args: args{
				s: "",
			},
			want:    0,
			want1:   0,
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := parseContentRange(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContenTRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseContenTRange() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseContenTRange() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseContenTRange() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_createSizedFile(t *testing.T) {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "testfile")
	defer os.RemoveAll(filePath)

	type args struct {
		path string
		size int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "createSizedFile successfully",
			args: args{
				path: filePath,
				size: 1000,
			},
			wantErr: false,
		},

		{
			name: "createSizedFile error, file size is 0",
			args: args{
				path: filePath,
				size: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createSizedFile(tt.args.path, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("createSizedFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModelDir(t *testing.T) {
	rootDir := "testValidateModelDir"
	modelDir := path.Join(rootDir, "testModel")
	testFile := "testValidateModelFile"
	os.MkdirAll(modelDir, 0755)
	os.Create(path.Join(modelDir, testFile))
	defer os.RemoveAll(rootDir)

	type args struct {
		dir   string
		model *Model
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "validateModelDir successfully",
			args: args{
				dir: rootDir,
				model: &Model{
					Format:    "SavedModel",
					FrameWork: "FrameWork",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModelDir(tt.args.dir, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("validateModelDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteORMBFile(t *testing.T) {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "test.yaml")
	defer os.RemoveAll(filePath)

	type args struct {
		filePath string
		model    *Model
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Write yaml file successfully",
			args: args{
				filePath: filePath,
				model: &Model{
					Format:    "SavedModel",
					FrameWork: "FrameWork",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeORMBFile(tt.args.filePath, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("WriteORMBFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_downloadModelFromHarbor(t *testing.T) {
	type args struct {
		client      ormb.Interface
		tenant      string
		user        string
		modelName   string
		versionName string
		model       *Model
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "downloadModelFromHarbor successfully",
			args: args{
				client: ormbClient,
				tenant: "system-tenant",
				user:   "admin",
				model: &Model{
					ModelName:   "test",
					VersionName: "v1",
					Format:      "savedmodel",
					FrameWork:   "tensorflow",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(path.Join(modelTmpDir, tt.args.tenant, tt.args.user, tt.args.modelName, tt.args.versionName))

			_, err := downloadModelFromHarbor(tt.args.client, tt.args.tenant, tt.args.user, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("downloadModelFromHarbor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
