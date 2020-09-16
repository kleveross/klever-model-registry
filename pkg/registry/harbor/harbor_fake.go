package harbor

import "net/http"

type fakeProxy struct {
}

func NewFakeProxy() ProxyClient {
	return &fakeProxy{}
}

func (p *fakeProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

func (p *fakeProxy) createModelJob(path string, byteManifests []byte) error {
	return nil
}

func (p *fakeProxy) ListArtifacts(project, repo string) ([]Artifact, error) {
	var testArtifacts []Artifact
	if project == "release" && repo == "tensorrt" {
		testArtifacts = append(testArtifacts, Artifact{
			Tags: []*Tag{
				{
					Name: "v1",
				},
			},
			ExtraAttrs: map[string]interface{}{
				"Author": "Klever",
				"Format": "TensorRT",
			},
		})
	}
	if project == "release" && repo == "savedmodel" {
		testArtifacts = append(testArtifacts, Artifact{
			Tags: []*Tag{
				{
					Name: "v1",
				},
			},
			ExtraAttrs: map[string]interface{}{
				"Author":    "Klever",
				"Format":    "SavedModel",
				"Framework": "TensorFlow",
			},
		})
	}
	return testArtifacts, nil
}
