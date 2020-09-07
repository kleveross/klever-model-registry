package harbor

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/caicloud/nirvana/log"
	"github.com/spf13/viper"
)

const (
	envModelRestirtyExternalAddress = "EXTERNAL_ADDRESS"
)

// Proxy is the proxy to Harbor core service.
type Proxy struct {
	Domain   string
	Username string
	Password string
}

// NewProxy creates a new proxy.
func NewProxy(domain, username, password string) *Proxy {
	return &Proxy{
		Domain:   domain,
		Username: username,
		Password: password,
	}
}

func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.Director = func(req *http.Request) {
		req.SetBasicAuth(p.Username, p.Password)

		req.Host = p.Domain
		req.URL.Host = p.Domain
		req.URL.Scheme = "http"
		req.URL.Path = r.URL.Path
		// If not add "/" suffix, harbor will match error router.
		if r.URL.Path == "/v2" || strings.HasSuffix(req.URL.Path, "/uploads") {
			req.URL.Path = r.URL.Path + "/"
		}
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		// Redirect's address is internal address in k8s, so MUST set model-registry's external address.
		if location, ok := resp.Header["Location"]; ok {
			resp.Header["Location"] = []string{strings.ReplaceAll(location[0], p.Domain,
				viper.GetString(envModelRestirtyExternalAddress))}
		}
		// It is to solve https://github.com/kleveross/klever-model-registry/issues/104
		resp.Header.Set("content-type", "application/json")
		return nil
	}

	var bodyBytes []byte
	var err error
	if strings.Contains(r.URL.Path, "manifests") {
		if r.Body != nil {
			bodyBytes, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("read request body err: %v", err)
			} else {
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
	}

	proxy.ServeHTTP(w, r)

	if strings.Contains(r.URL.Path, "manifests") && len(bodyBytes) > 0 {
		err := p.createModelJob(r.URL.Path, bodyBytes)
		if err != nil {
			log.Errorf("create modeljob error when push model, err: %v", err)
		}
	}
}
