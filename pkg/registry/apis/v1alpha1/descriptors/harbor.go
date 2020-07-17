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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/spf13/viper"

	"github.com/kleveross/klever-model-registry/pkg/common"
)

func HarborProxy(w http.ResponseWriter, r *http.Request) {
	ormbDomain := viper.GetString(common.ORMBDomainEnvKey)
	ormbUserName := viper.GetString(common.ORMBUsernameEnvkey)
	ormbPassword := viper.GetString(common.ORMBPasswordEnvKey)
	authBytes := []byte(fmt.Sprintf("%v:%v", ormbUserName, ormbPassword))
	auth := "basic " + base64.StdEncoding.EncodeToString(authBytes)

	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.Director = func(req *http.Request) {
		req.Header = r.Header
		req.Header["Authorization"] = []string{auth}
		req.Host = ormbDomain
		req.URL.Host = ormbDomain
		req.URL.Scheme = "http"
		req.URL.Path = r.URL.Path
	}

	proxy.ServeHTTP(w, r)
}
