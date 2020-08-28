package serving

import (
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	seldonv1client "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
)

type ServingController struct {
	seldonClient seldonv1client.Interface
}

func New(seldClient seldonv1client.Interface) *ServingController {
	return &ServingController{
		seldonClient: seldClient,
	}
}

func (s ServingController) Create(namespace string, sdep *seldonv1.SeldonDeployment) error {
	Compose(sdep)

	_, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Create(sdep)
	if err != nil {
		return errors.RenderError(err)
	}

	return nil
}

func (s ServingController) Get(namespace, sdepID string) (*seldonv1.SeldonDeployment, error) {
	sdep, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Get(sdepID, metav1.GetOptions{})
	if err != nil {
		return nil, errors.RenderError(err)
	}
	return sdep, nil
}

func (s ServingController) List(namespace string) (*seldonv1.SeldonDeploymentList, error) {
	sdeps, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.RenderError(err)
	}
	return sdeps, nil
}

func (s ServingController) Update(namespace string, sdep *seldonv1.SeldonDeployment) error {
	_, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Update(sdep)
	if err != nil {
		return errors.RenderError(err)
	}

	return nil
}

func (s ServingController) Delete(namespace, sdepID string) error {
	err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Delete(sdepID, &metav1.DeleteOptions{})
	if err != nil {
		return errors.RenderError(err)
	}
	return nil
}
