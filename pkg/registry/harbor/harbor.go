package harbor

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/spf13/viper"

	"github.com/kleveross/klever-model-registry/pkg/common"
)

// Proxy is the reverse proxy to Harbor.
func Proxy(w http.ResponseWriter, r *http.Request) {
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
