package models

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/caicloud/nirvana/log"
	ormbmodel "github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/ormb"
	"gopkg.in/yaml.v2"

	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/util"
)

const (
	// uploadMaxSize is the max size of upload file
	uploadMaxSize = 100 << 20
)

// errNoContentRange is the error which http header has not Content-Range key.
var errNoContentRange = errors.New("no Content-Range")

// errSeeker is returned by ServeContent's sizeFunc when the content
// doesn't seek properly. The underlying Seeker's error text isn't
// included in the sizeFunc reply so it's not sent over HTTP to end
// users.
var errSeeker = errors.New("seeker can't seek")

// errNoOverlap is returned by serveContent's parseRange if first-byte-pos of
// all of the byte-range-spec values is greater than the content size.
var errNoOverlap = errors.New("invalid range: failed to overlap")

func validateFileSize(w http.ResponseWriter, r *http.Request) error {
	if contentLength := r.Header.Get("Content-Length"); contentLength != "" {
		v, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return fmt.Errorf("http: invalid Content-Length of %v", contentLength)
		}
		if v > uploadMaxSize {
			return fmt.Errorf("upload file fail: file to large %v", v)
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, uploadMaxSize)
	if err := r.ParseMultipartForm(uploadMaxSize); err != nil {
		return fmt.Errorf("failed parse form file: %v", err)
	}
	return nil
}

type ChunkInfo struct {
	Path      string
	FileName  string
	Content   io.ReadCloser
	TotalSize int64
	PartFrom  int64
	PartTo    int64
}

func parseChunk(r *http.Request) (*ChunkInfo, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed parse file form request: %v", err)
	}

	totalSize, partFrom, partTo, err := parseContentRange(r.Header.Get("Content-Range"))
	if err != nil && err != errNoContentRange {
		return nil, fmt.Errorf("failed to parse Content-Range: %v", err)
	}

	if err == errNoContentRange {
		totalSize, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, errSeeker
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, errSeeker
		}
		partFrom, partTo = 0, totalSize-1
	}

	return &ChunkInfo{
		FileName:  fileHeader.Filename,
		Content:   file,
		TotalSize: totalSize,
		PartFrom:  partFrom,
		PartTo:    partTo,
	}, nil
}

func parseContentRange(s string) (int64, int64, int64, error) {
	if s == "" {
		return 0, 0, 0, errNoContentRange
	}

	s = strings.Replace(s, "bytes ", "", -1)
	fromTo := strings.Split(s, "/")[0]
	totalSize, err := strconv.ParseInt(strings.Split(s, "/")[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	splitted := strings.Split(fromTo, "-")

	partFrom, err := strconv.ParseInt(splitted[0], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	partTo, err := strconv.ParseInt(splitted[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return totalSize, partFrom, partTo, nil
}

func createSizedFile(path string, size int64) error {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := fd.Close(); cerr != nil {
			log.Errorf("[createSizedFile] close file err: %s", cerr.Error())
		}
	}()

	if size <= 0 {
		return fmt.Errorf("file size can not be zero")
	}

	_, err = fd.Seek(size-1, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = fd.Write([]byte{0})
	return err
}

func downloadModelFromHarbor(client ormb.Interface, tenant, user string, model *Model) (string, error) {
	filePath := path.Join(modelTmpDir, tenant, user, model.ModelName, model.VersionName)
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(filePath)

	modelRef := fmt.Sprintf("%v/%v/%v:%v", common.ORMBDomain, model.ModelName,
		model.ModelName, model.VersionName)
	err = client.Pull(modelRef)
	if err != nil {
		return "", err
	}
	defer func() {
		err := client.Remove(modelRef)
		if err != nil {
			log.Warningf("Remove model err: %v", err)
		}
	}()

	err = client.Export(modelRef, filePath)
	if err != nil {
		return "", err
	}

	zipFileName := filePath + ".zip"
	err = util.Archive(filePath, zipFileName)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(zipFileName)

	return zipFileName, nil
}

func uploadModelToHarbor(client ormb.Interface, zipFile string, model *Model) error {
	deCompressDir := strings.Trim(zipFile, ".zip")

	defer func() {
		// No matter success or failure, MUST delete unzrchive dir.
		err := os.RemoveAll(deCompressDir)
		if err != nil {
			log.Warningf("Remove %v err: %v", deCompressDir, err)
		}
	}()

	err := util.Unarchive(zipFile, deCompressDir)
	if err != nil {
		return err
	}

	err = validateModelDir(deCompressDir, model)
	if err != nil {
		return err
	}

	modelRef := fmt.Sprintf("%v/%v/%v:%v", common.ORMBDomain, model.ModelName,
		model.ModelName, model.VersionName)
	err = client.Save(deCompressDir, modelRef)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Remove(modelRef)
		if err != nil {
			log.Warningf("Remove model err: %v", err)
		}
	}()

	err = client.Push(modelRef)
	if err != nil {
		return err
	}

	return nil
}

func validateModelDir(dirPath string, model *Model) error {
	rootDir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	if len(rootDir) != 1 || !rootDir[0].IsDir() {
		return fmt.Errorf("check model dir err， need one dir in dir %s ", dirPath)
	}
	// ORMB dir structure as follow
	// -| model
	//      -| modelFile
	//    ormbfile.yaml
	ormbModelDir := path.Join(dirPath, "model")
	err = os.MkdirAll(ormbModelDir, 0755)
	if err != nil {
		return err
	}
	fileListPath := path.Join(dirPath, rootDir[0].Name())
	fileList, err := ioutil.ReadDir(fileListPath)
	if err != nil {
		return err
	}
	for _, file := range fileList {
		err = util.ExecOSCommand("mv",
			[]string{
				path.Join(fileListPath, file.Name()),
				ormbModelDir})

		if err != nil {
			return err
		}
	}
	err = util.ExecOSCommand("rm", []string{"-rf", path.Join(dirPath, rootDir[0].Name())})
	if err != nil {
		return err
	}

	err = writeORMBFile(path.Join(dirPath, "ormbfile.yaml"), model)
	if err != nil {
		return err
	}

	return nil
}

func writeORMBFile(filePath string, model *Model) error {
	metadata := ormbmodel.Metadata{
		Author:      "",
		Description: model.Description,
		Format:      model.Format,
		Framework:   model.Format,
		Signature: ormbmodel.Signature{
			Inputs:  model.Inputs,
			Outputs: model.Outputs,
		},
	}
	data, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
