/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package descriptors

import (
	"net/http"

	"github.com/caicloud/nirvana/definition"

	"github.com/kleveross/klever-model-registry/pkg/registry/harbor"
)

// HarborAPIDescriptor contain horbor /api/* descriptors
func HarborAPIDescriptor() definition.Descriptor {
	return definition.Descriptor{
		Path:        "/api/{path:*}",
		Description: "It contains all harbor /api/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEJSON},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: http.HandlerFunc(harbor.Proxy),
			},
		},
	}
}

// HarborServiceDescriptor contain horbor /service/* descriptors
func HarborServiceDescriptor() definition.Descriptor {
	return definition.Descriptor{
		Path:        "/service/{path:*}",
		Description: "It contains all api in /service/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEJSON},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: http.HandlerFunc(harbor.Proxy),
			},
		},
	}
}

// HarborCDescriptor contain horbor /c/* descriptors
func HarborCDescriptor() definition.Descriptor {
	return definition.Descriptor{
		Path:        "/c/{path:*}",
		Description: "It contains all api in /c/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEJSON},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: http.HandlerFunc(harbor.Proxy),
			},
		},
	}
}

// HarborV2Descriptor contain horbor /v2/* descriptors
func HarborV2Descriptor() definition.Descriptor {
	return definition.Descriptor{
		Path:        "/v2/{path:*}",
		Description: "It contains all api in /v2/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEJSON},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: http.HandlerFunc(harbor.Proxy),
			},
		},
	}
}
