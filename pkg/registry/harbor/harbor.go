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

// ProxyClient is the proxy client to Harbor core service.
type ProxyClient interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	createModelJob(path string, byteManifests []byte) error
	ListArtifacts(project, repo string) ([]Artifact, error)
}

// proxy is the proxy to Harbor core service.
type proxy struct {
	Domain   string
	Username string
	Password string
}

func NewProxy(domain, username, password string) ProxyClient {
	return &proxy{
		Domain:   domain,
		Username: username,
		Password: password,
	}
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		// There are errors for cookies, eg: there are error like `"message": "CSRF token invalid"`.
		// So we set no cookies to avoid it.
		resp.Header["Set-Cookie"] = []string{}

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
