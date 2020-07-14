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
package harbor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/astaxie/beego/context"
	"github.com/spf13/viper"

	"github.com/caicloud/temp-model-registry/pkg/common"
	"github.com/caicloud/temp-model-registry/pkg/registry/types"
)

func webhookAPI(c *context.Context) {
	payload := types.Payload{}
	if err := json.Unmarshal(c.Input.RequestBody, &payload); err != nil {
		panic(err)
	}
	// TODO(simon): Create ModelJob to extract or convert

	c.ResponseWriter.Write([]byte("OK"))
}

func harborProxyAPIs(c *context.Context) {
	ormbDomain := viper.GetString(common.ORMBDomainEnvKey)
	ormbUserName := viper.GetString(common.ORMBUsernameEnvkey)
	ormbPassword := viper.GetString(common.ORMBPasswordEnvKey)
	authBytes := []byte(fmt.Sprintf("%v:%v", ormbUserName, ormbPassword))
	auth := "basic " + base64.StdEncoding.EncodeToString(authBytes)

	proxy := httputil.NewSingleHostReverseProxy(c.Request.URL)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Header["Authorization"] = []string{auth}
		req.Host = ormbDomain
		req.URL.Host = ormbDomain
		req.URL.Scheme = "http"
		req.URL.Path = c.Request.URL.Path
	}

	proxy.ServeHTTP(c.ResponseWriter, c.Request)
}
