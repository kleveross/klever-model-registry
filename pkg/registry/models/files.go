package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/caicloud/nirvana/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/modeljob"
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

// UploadFile uploads the model file to harbor
func UploadFile(ctx context.Context, tenant, user, projectName, modelName, versionName string) error {
	var model Model
	modelContent := util.GetFormValueFromRequest(ctx, "model")
	err := json.Unmarshal([]byte(modelContent), &model)
	if err != nil {
		return errors.RenderBadRequestError(err)
	}
	model.ProjectName = projectName
	model.ModelName = modelName
	model.VersionName = versionName

	modelDir := path.Join(modelTmpDir, tenant, user, modelName)
	err = os.MkdirAll(modelDir, 0755)
	if err != nil {
		return errors.RenderInternalServerError(err)
	}

	request := util.GetRequestFromContext(ctx)
	responseWriter := util.GetResponseFromContext(ctx)
	err = validateFileSize(responseWriter, request)
	if err != nil {
		log.Errorf("Failed to validate the file size: %v", err)
		return errors.RenderBadRequestError(err)
	}

	chunkInfo, err := parseChunk(request)
	if err != nil {
		log.Errorf("Failed to parse the chunk for the form key `file`: %v", err)
		return errors.RenderBadRequestError(err)
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
			log.Errorf("Failed to create the file: %v", err)
			return errors.RenderInternalServerError(err)
		}
	}

	if chunkInfo.TotalSize != 0 {
		newFile, err := os.OpenFile(zipFileName, os.O_WRONLY, 0)
		if err != nil {
			return errors.RenderInternalServerError(err)
		}
		defer func() {
			if cerr := newFile.Close(); cerr != nil {
				log.Errorf("Failed to close: %v", cerr)
			}
		}()
		_, err = newFile.Seek(chunkInfo.PartFrom, io.SeekStart)
		if err != nil {
			log.Errorf("Failed to seek: %v", err)
			return errors.RenderInternalServerError(err)
		}
		_, err = io.Copy(newFile, chunkInfo.Content)
		if err != nil {
			log.Errorf("Failed to copy the content to the new created file: %v", err)
			return errors.RenderInternalServerError(err)
		}
	}

	if chunkInfo.TotalSize-1 == chunkInfo.PartTo {
		err = uploadModelToHarbor(client.GetORMBClient(), zipFileName, &model)
		if err != nil {
			log.Errorf("Failed to update the model to harbor: %v", err)
			return errors.RenderInternalServerError(err)
		}

		if modeljob.IsExtractModel(model.Format) {
			modeljobObj := modeljob.GenerateExtractionModelJob(common.ORMBDomain, projectName, modelName, versionName, model.Format)
			_, err = client.GetKubeKleverOssClient().KleverossV1alpha1().ModelJobs("default").Create(context.Background(), modeljobObj, metav1.CreateOptions{})
			if err != nil {
				return errors.RenderInternalServerError(err)
			}
		}
	}

	return nil
}

// DownloadFile downloads the model file from harbor
func DownloadFile(ctx context.Context, tenant, user, projectName, modelName, versionName string, model *Model) error {
	model.ProjectName = projectName
	model.ModelName = modelName
	model.VersionName = versionName

	zipFileName, err := downloadModelFromHarbor(client.GetORMBClient(), tenant, user, model)
	if err != nil {
		return errors.RenderInternalServerError(err)
	}

	file, err := os.Open(zipFileName)
	if err != nil {
		return errors.RenderInternalServerError(err)
	}
	defer file.Close()

	responseWriter := util.GetResponseFromContext(ctx)
	responseWriter.Header().Set("Content-Type", "application/octet-stream")
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.%s\"", versionName, "zip"))
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		return errors.RenderInternalServerError(err)
	}

	return nil
}
