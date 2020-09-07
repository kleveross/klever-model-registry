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
	"github.com/caicloud/nirvana/definition"

	"github.com/kleveross/klever-model-registry/pkg/registry/harbor"
)

var harborController *harbor.Proxy

// InitHarborController handles the harbor proxy
func InitHarborController(domain, username, password string) {
	harborController = harbor.NewProxy(domain, username, password)
}

// HarborAPIPrefixDescriptor contain horbor /api/* descriptors
func HarborAPIPrefixDescriptor(domain, username, password string) definition.Descriptor {
	return definition.Descriptor{
		Path:        "/api/{path:*}",
		Description: "It contains all harbor /api/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEAll},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: harbor.NewProxy(domain, username, password),
			},
		},
	}
}

// HarborServicePrefixDescriptor contain horbor /service/* descriptors
func HarborServicePrefixDescriptor(domain, username, password string) definition.Descriptor {
	return definition.Descriptor{
		Path:        "/service/{path:*}",
		Description: "It contains all api in /service/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEAll},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: harbor.NewProxy(domain, username, password),
			},
		},
	}
}

// HarborCPrefixDescriptor contain horbor /c/* descriptors
func HarborCPrefixDescriptor(domain, username, password string) definition.Descriptor {
	return definition.Descriptor{
		Path:        "/c/{path:*}",
		Description: "It contains all api in /c/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEAll},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: harbor.NewProxy(domain, username, password),
			},
		},
	}
}

// HarborV2PrefixDescriptor contain horbor /v2/* descriptors
func HarborV2PrefixDescriptor(domain, username, password string) definition.Descriptor {
	return definition.Descriptor{
		Path:        "/v2/{path:*}",
		Description: "It contains all api in /v2/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEAll},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: harbor.NewProxy(domain, username, password),
			},
		},
	}
}

func HarborV2Descriptor(domain, username, password string) definition.Descriptor {
	return definition.Descriptor{
		Path:        "/v2/",
		Description: "It contains all api in /v2/*",
		Consumes:    []string{definition.MIMEAll},
		Produces:    []string{definition.MIMEAll},
		Definitions: []definition.Definition{
			{
				Method:  definition.Any,
				Handler: harbor.NewProxy(domain, username, password),
			},
		},
	}
}
