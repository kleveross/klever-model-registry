package modeljob_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
)

var _ = Describe("Modeljob API", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1

	modeljobObj := &modeljobsv1alpha1.ModelJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "modelJobName",
		},
	}

	It("Should manager modeljob successfully", func() {
		// Create modeljob
		modeljobCreated, err := modeljobController.Create("default", modeljobObj)
		Expect(err).To(BeNil())

		// Get modeljob
		By("Expecting get modeljob successfully")
		Eventually(func() error {
			_, err := modeljobController.Get(modeljobCreated.Namespace, modeljobCreated.Name)
			return err
		}, timeout, interval).Should(Succeed())

		// List modeljob
		modeljobeListed, err := modeljobController.List(modeljobCreated.Namespace, "", &paging.ListOption{
			Start: 0,
		})
		Expect(err).To(BeNil())
		Expect(modeljobeListed.ListMeta.TotalItems).To(Equal(1))

		// Delete modeljob
		err = modeljobController.Delete(modeljobCreated.Namespace, modeljobCreated.Name)
		Expect(err).To(BeNil())
	})
})
