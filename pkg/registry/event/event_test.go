package event_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
)

var _ = Describe("Event", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1

	modeljobObj := &modeljobsv1alpha1.ModelJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "modelJobName",
		},
	}

	It("Should manager modeljob events successfully", func() {
		// Create modeljob
		modeljobCreated, err := modeljobController.Create("default", modeljobObj)
		Expect(err).To(BeNil())
		limit := 10

		// Get modeljob events
		By("Expecting get event successfully")
		Eventually(func() error {
			_, err := eventController.GetModelJobEvents("default", modeljobCreated.Name, &paging.ListOption{
				Start: 0,
				Limit: &limit,
			})
			return err
		}, timeout, interval).Should(Succeed())

		// Get modeljob events, have error modeljob name
		_, err = eventController.GetModelJobEvents("default", "nonModelJobName", &paging.ListOption{
			Start: 0,
			Limit: &limit,
		})
		Expect(err).To(HaveOccurred())
	})
})
