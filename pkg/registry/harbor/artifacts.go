package harbor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/caicloud/nirvana/log"

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/modeljob"
)

// Artifact is copy from https://github.com/goharbor/harbor/blob/master/src/pkg/artifact/model.go#L31-L47
// can not reference `Artifact` struct directly, since it cant not find
// https://github.com/goharbor/harbor/blob/master/src/pkg/artifact/model.go#L24 package.
type Artifact struct {
	ID                int64                  `json:"id"`
	Type              string                 `json:"type"`                // image, chart, etc
	MediaType         string                 `json:"media_type"`          // the media type of artifact. Mostly, it's the value of `manifest.config.mediatype`
	ManifestMediaType string                 `json:"manifest_media_type"` // the media type of manifest/index
	ProjectID         int64                  `json:"project_id"`
	RepositoryID      int64                  `json:"repository_id"`
	RepositoryName    string                 `json:"repository_name"`
	Digest            string                 `json:"digest"`
	Size              int64                  `json:"size"`
	Icon              string                 `json:"icon"`
	PushTime          time.Time              `json:"push_time"`
	PullTime          time.Time              `json:"pull_time"`
	ExtraAttrs        map[string]interface{} `json:"extra_attrs"` // only contains the simple attributes specific for the different artifact type, most of them should come from the config layer
	Annotations       map[string]string      `json:"annotations"`
	Tags              []*Tag                 `json:"tags"` // the list of tags that attached to the artifact
}

// Tag is copy from https://github.com/goharbor/harbor/blob/master/src/pkg/tag/model/tag/model.go
type Tag struct {
	ID           int64     `json:"id"`
	RepositoryID int64     `json:"repository_id"` // tags are the resources of repository, one repository only contains one same name tag
	ArtifactID   int64     `json:"artifact_id"`   // the artifact ID that the tag attaches to, it changes when pushing a same name but different digest artifact
	Name         string    `json:"name"`
	PushTime     time.Time `json:"push_time"`
	PullTime     time.Time `json:"pull_time"`
}

func (p *Proxy) createModelJob(path string, byteManifests []byte) error {
	path = strings.Trim(path, "/")
	pathSlice := strings.Split(path, "/")
	projectName := pathSlice[1]
	modelName := pathSlice[2]
	versionName := pathSlice[4]

	artis, err := p.ListArtifacts(projectName, modelName)
	if err != nil {
		return err
	}

	var found *Artifact
	for artiIndex, arti := range artis {
		for _, tag := range arti.Tags {
			if tag.Name == versionName {
				found = &artis[artiIndex]
				break
			}
		}
		if found != nil {
			break
		}
	}

	if found == nil {
		log.Info("not found match artifacts")
		return nil
	}
	if format, ok := found.ExtraAttrs["format"]; ok {
		modeljobObj := modeljob.GenerateExtractionModelJob(p.Domain, projectName, modelName, versionName, format.(string))
		_, err := client.GetKubeKleverOssClient().KleverossV1alpha1().ModelJobs("default").Create(modeljobObj)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Proxy) ListArtifacts(project, repo string) ([]Artifact, error) {
	url := fmt.Sprintf("http://%v/api/v2.0/projects/%v/repositories/%v/artifacts",
		p.Domain, project, repo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.Username, p.Password)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var artis []Artifact
	err = json.Unmarshal(bodyBytes, &artis)
	if err != nil {
		return nil, err
	}

	return artis, nil
}
