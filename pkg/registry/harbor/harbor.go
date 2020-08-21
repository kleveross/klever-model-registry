package harbor

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
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
	authBytes := []byte(fmt.Sprintf("%v:%v", p.Username, p.Password))
	auth := "basic " + base64.StdEncoding.EncodeToString(authBytes)

	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.Director = func(req *http.Request) {
		req.Header = r.Header
		req.Header["Authorization"] = []string{auth}
		req.Host = p.Domain
		req.URL.Host = p.Domain
		req.URL.Scheme = "http"
		req.URL.Path = r.URL.Path
	}

	proxy.ServeHTTP(w, r)
}
