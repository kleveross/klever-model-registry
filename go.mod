module github.com/kleveross/klever-model-registry

go 1.14

require (
	github.com/caicloud/nirvana v0.2.10
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/gavv/httpexpect/v2 v2.1.0
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.4.4
	github.com/golang/snappy v0.0.1 // indirect
	github.com/kleveross/ormb v0.0.5
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/prometheus/client_golang v1.7.0 // indirect
	github.com/seldonio/seldon-core/operator v0.0.0-20200714192520-b8cf277bc2ad
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.7.0
	github.com/ulikunitz/xz v0.5.7 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v0.17.9
	sigs.k8s.io/controller-runtime v0.5.8
)

replace github.com/seldonio/seldon-core/operator => github.com/kleveross/seldon-core/operator v0.0.0-20200917064040-3e7d8727897f
