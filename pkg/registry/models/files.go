package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/caicloud/nirvana/log"

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/util"
)

var (
	modelTmpDir string
)

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	modelTmpDir = path.Join(workDir, "modelx")
}

func UploadFile(ctx context.Context, tenant, user, modelName, versionName string) error {
	var model Model
	modelContent := util.GetFormValueFromRequest(ctx, "model")
	err := json.Unmarshal([]byte(modelContent), &model)
	if err != nil {
		return err
	}

	modelDir := path.Join(modelTmpDir, tenant, user, modelName)
	err = os.MkdirAll(modelDir, 0755)
	if err != nil {
		return err
	}

	request := util.GetRequstFromContext(ctx)
	responseWriter := util.GetResponseFromContext(ctx)
	err = validateFileSize(responseWriter, request)
	if err != nil {
		return err
	}

	chunkInfo, err := parseChunk(request)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := chunkInfo.Content.Close(); cerr != nil {
			log.Errorf("chunInfo.Content close err: %v", cerr.Error())
		}
	}()
	zipFileName := path.Join(modelDir, versionName+".zip")

	if chunkInfo.PartFrom == 0 {
		err = createSizedFile(zipFileName, chunkInfo.TotalSize)
		if err != nil {
			return err
		}
	}

	if chunkInfo.TotalSize != 0 {
		newFile, err := os.OpenFile(zipFileName, os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := newFile.Close(); cerr != nil {
				log.Errorf("failed to close: %v", cerr)
			}
		}()
		_, err = newFile.Seek(chunkInfo.PartFrom, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = io.Copy(newFile, chunkInfo.Content)
		if err != nil {
			return err
		}
	}

	if chunkInfo.TotalSize-1 == chunkInfo.PartTo {
		err = uploadModelToHarbor(client.ORMBClient, zipFileName, &model)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadFile(ctx context.Context, tenant, user, modelName, versionName string, model *Model) error {
	zipFileName, err := downloadModelFromHarbor(client.ORMBClient, tenant, user, modelName, versionName, model)
	if err != nil {
		return err
	}

	file, err := os.Open(zipFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	responseWriter := util.GetResponseFromContext(ctx)
	responseWriter.Header().Set("Content-Type", "application/octet-stream")
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.%s\"", versionName, "zip"))
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		return err
	}

	return nil
}
