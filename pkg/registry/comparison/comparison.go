package comparison

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/caicloud/nirvana/log"
	ormbmodel "github.com/kleveross/ormb/pkg/model"

	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/registry/harbor"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
	"github.com/kleveross/klever-model-registry/pkg/util"
)

// Generator list models' metadata and compare
func Generator(ctx context.Context, models Comparison, opt *paging.ListOption) (*ORMBModelList, error) {
	proxy := harbor.NewProxy(common.ORMBDomain, common.ORMBUserName, common.ORMBPassword)
	metaList, err := composeComparison(models.Models, proxy)
	if err != nil {
		return nil, err
	}
	return toORMBModelList(metaList, opt), nil
}

// toModelJobList is convert to ModelJobList struct.
func toORMBModelList(items []*ormbmodel.Model, opt *paging.ListOption) *ORMBModelList {
	datas := paging.Page(items, opt)
	modelList := &ORMBModelList{
		ListMeta: paging.ListMeta{
			TotalItems: datas.TotalItems,
		},
		Items: []*ormbmodel.Model{},
	}

	for _, d := range datas.Items {
		modelList.Items = append(modelList.Items, d.(*ormbmodel.Model))
	}
	return modelList
}

func composeComparison(models []ComparisonModel, proxy harbor.ProxyClient) ([]*ormbmodel.Model, error) {
	metaList := make([]*ormbmodel.Model, 0)
	for _, model := range models {
		artifacts, err := proxy.ListArtifacts(model.Project, model.Name)
		if err != nil {
			return nil, err
		}
		for _, artifact := range artifacts {
			for _, tag := range artifact.Tags {
				if tag.Name == model.Tag {
					manifest, err := json.Marshal(artifact.ExtraAttrs)
					if err != nil {
						return nil, err
					}
					var meta ormbmodel.Metadata
					if err := json.Unmarshal(manifest, &meta); err != nil {
						return nil, err
					}
					metaList = append(metaList, &ormbmodel.Model{
						Path:     fmt.Sprintf("%s/%s:%s", model.Project, model.Name, model.Tag),
						Metadata: &meta,
					})
					break
				}
			}
		}
	}

	return metaList, nil
}

// DownloadCSVFile downloads the comparison csv file
func DownloadCSVFile(ctx context.Context, models Comparison) error {
	proxy := harbor.NewProxy(common.ORMBDomain, common.ORMBUserName, common.ORMBPassword)
	metaList, err := composeComparison(models.Models, proxy)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%d.csv", time.Now().Unix())
	err = generateCSVFile(ctx, fileName, metaList)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(fileName); err != nil {
			log.Error("remove file is failed, err: %s", err.Error())
		}
	}()

	err = responseCSVFile(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}

func generateCSVFile(ctx context.Context, fileName string, metas []*ormbmodel.Model) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("open file is failed, err: %s", err.Error())
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error("close file is failed, err: %s", err.Error())
		}
	}()

	// Write UTF-8 BOM mainly for unidentifiable Chinese code,
	// see https://stackoverflow.com/questions/2223882/whats-the-difference-between-utf-8-and-utf-8-without-bom
	_, err = file.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return fmt.Errorf("write utf-8 bom err: %s", err.Error())
	}
	csvWrite := csv.NewWriter(file)

	content, err := composeCSVFileContent(metas)
	if err != nil {
		return err
	}
	err = csvWrite.WriteAll(content)
	if err != nil {
		return err
	}
	csvWrite.Flush()

	return nil
}

func responseCSVFile(ctx context.Context, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error("close file is failed, err: %s", err.Error())
		}
	}()

	responseWriter := util.GetResponseFromContext(ctx)
	responseWriter.Header().Set("Content-Type", "application/octet-stream")
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		return err
	}

	return nil
}

func composeCSVFileContent(metas []*ormbmodel.Model) ([][]string, error) {
	fileContents := [][]string{
		{"Basic Info"},
		{"Model Name"},
		{"Model Source"},
		{"Model Framework"},
		{"Model Format"},
		{"Model Inputs"},
		{"Model Outputs"},
	}

	for index, content := range fileContents {
		for _, meta := range metas {
			switch content[0] {
			case "Basic Info":
				continue
			case "Model Name":
				content = append(content, meta.Path)
			case "Model Source":
				content = append(content, meta.Metadata.Author)
			case "Model Framework":
				content = append(content, meta.Metadata.Framework)
			case "Model Format":
				content = append(content, meta.Metadata.Format)
			case "Model Inputs":
				jsonString := "-"
				if meta.Metadata.Signature != nil && meta.Metadata.Signature.Inputs != nil {
					jsonString = composeJSONString(meta.Metadata.Signature.Inputs)
				}
				content = append(content, jsonString)
			case "Model Outputs":
				jsonString := "-"
				if meta.Metadata.Signature != nil && meta.Metadata.Signature.Outputs != nil {
					jsonString = composeJSONString(meta.Metadata.Signature.Outputs)
				}
				content = append(content, jsonString)
			}
			fileContents[index] = content
		}
	}

	return fileContents, nil
}

func composeJSONString(obj interface{}) string {
	bytes, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		log.Errorf("Compose JSON string failed: %v", err.Error())
		return ""
	}
	return string(bytes)
}
