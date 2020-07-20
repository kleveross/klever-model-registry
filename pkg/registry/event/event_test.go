package event_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	modeljobfake "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned/fake"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/event"
	"github.com/kleveross/klever-model-registry/pkg/registry/modeljob"
)

var _ = Describe("Event", func() {
	client.KubeMainClient = testclient.NewSimpleClientset()
	client.KubeModelJobClient = modeljobfake.NewSimpleClientset()

	modeljobObj := &modeljobsv1alpha1.ModelJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "modelJobName",
		},
	}

	// Create modeljob
	modeljobCreated, err := modeljob.Create(modeljobObj)
	Expect(err).To(BeNil())

	// Get modeljob events
	_, err = event.GetModelJobEvents("default", modeljobCreated.Name)
	Expect(err).To(BeNil())

	// Get modeljob events, have error modeljob name
	_, err = event.GetModelJobEvents("default", "nonModelJobName")
	Expect(err).To(HaveOccurred())
})
