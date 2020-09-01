package integration

import "os"

var (
	// ModelRegistryHost is the host of the model registry.
	ModelRegistryHost = "http://127.0.0.1:30002"
	// EnvVarModelRegistryHost is the model registry host environment variable.
	EnvVarModelRegistryHost = "MODEL_REGISTRY_HOST"
)

func init() {
	if val := os.Getenv(EnvVarModelRegistryHost); val != "" {
		ModelRegistryHost = val
	}
}
