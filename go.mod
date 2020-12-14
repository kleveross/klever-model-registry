module github.com/kleveross/klever-model-registry

go 1.14

require (
	github.com/caicloud/nirvana v0.2.10
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/gavv/httpexpect/v2 v2.1.0
	github.com/go-logr/logr v0.3.0
	github.com/golang/mock v1.4.4
	github.com/kleveross/ormb v0.0.8
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/seldonio/seldon-core/operator v1.5.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.4
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/go-logr/logr => github.com/go-logr/logr v0.1.0
	github.com/seldonio/seldon-core/operator => github.com/kleveross/seldon-core/operator v0.0.0-20201214071233-b0fc794d2e91
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	k8s.io/client-go => k8s.io/client-go v0.19.3
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.1.0
)
