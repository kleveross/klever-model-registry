package modeljob_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	modeljobfake "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned/fake"
	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/modeljob"
)

var _ = Describe("Modeljob API", func() {
	client.KubeModelJobClient = modeljobfake.NewSimpleClientset()

	modeljobObj := &modeljobsv1alpha1.ModelJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "modelJobName",
		},
	}

	// Create modeljob
	modeljobCreated, err := modeljob.Create(common.DefaultModelJobNamespace, modeljobObj)
	Expect(err).To(BeNil())

	// Get modeljob
	modeljobGeted, err := modeljob.Get(modeljobCreated.Namespace, modeljobCreated.Name)
	Expect(err).To(BeNil())
	Expect(modeljobGeted.Name).To(Equal(modeljobCreated.Name))

	// List modeljob
	modeljobeListed, err := modeljob.List(modeljobCreated.Namespace)
	Expect(err).To(BeNil())
	Expect(len(modeljobeListed.Items)).To(Equal(1))

	// Delete modeljob
	err = modeljob.Delete(modeljobCreated.Namespace, modeljobCreated.Name)
	Expect(err).To(BeNil())
})
