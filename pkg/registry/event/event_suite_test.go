package event_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"

	modeljobfake "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned/fake"
	kleverossinformers "github.com/kleveross/klever-model-registry/pkg/clientset/informers/externalversions"
	"github.com/kleveross/klever-model-registry/pkg/registry/event"
	"github.com/kleveross/klever-model-registry/pkg/registry/modeljob"
)

var (
	eventController    *event.EventController
	modeljobController *modeljob.ModelJobController
	stopCh             = make(chan struct{})
)

func TestEvent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Suite")
}

var _ = BeforeSuite(func() {
	mainClient := k8sfake.NewSimpleClientset()
	kleverossClient := modeljobfake.NewSimpleClientset()
	factory := kleverossinformers.NewSharedInformerFactory(kleverossClient, 30*time.Second)

	go factory.Start(stopCh)

	modeljobInformer := factory.Kleveross().V1alpha1().ModelJobs()
	if !cache.WaitForCacheSync(stopCh, modeljobInformer.Informer().HasSynced) {
		panic(fmt.Errorf("failed to wait for modeljob synced"))
	}

	eventController = event.New(mainClient, kleverossClient, modeljobInformer)
	modeljobController = modeljob.New(kleverossClient, modeljobInformer)
})

var _ = AfterSuite(func() {
	stopCh <- struct{}{}
})
