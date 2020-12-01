package serving

import (
	"context"
	"sort"

	"github.com/caicloud/nirvana/log"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	seldonv1client "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
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
	if err := Compose(sdep); err != nil {
		log.Errorf("Failed to compose the Seldon Deployment: %v", err)
		return errors.RenderError(err)
	}

	_, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Create(context.TODO(), sdep, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create the Seldon Deployment: %v", err)
		return errors.RenderError(err)
	}

	return nil
}

func (s ServingController) Get(namespace, sdepID string) (*seldonv1.SeldonDeployment, error) {
	sdep, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Get(context.TODO(), sdepID, metav1.GetOptions{})
	if err != nil {
		return nil, errors.RenderError(err)
	}
	return sdep, nil
}

// ServingList is the response of List Interface.
type ServingList struct {
	ListMeta paging.ListMeta             `json:"metadata"`
	Items    []seldonv1.SeldonDeployment `json:"items"`
}

// toServingList is convert to ServingList struct
func toServingList(items []seldonv1.SeldonDeployment, opt *paging.ListOption) *ServingList {
	datas := paging.Page(items, opt)
	servingList := &ServingList{
		ListMeta: paging.ListMeta{
			TotalItems: datas.TotalItems,
		},
		Items: []seldonv1.SeldonDeployment{},
	}

	for _, d := range datas.Items {
		servingList.Items = append(servingList.Items, d.(seldonv1.SeldonDeployment))
	}
	return servingList
}

func (s ServingController) List(namespace string, opt *paging.ListOption) (*ServingList, error) {
	sdeps, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.RenderError(err)
	}

	sort.SliceStable(sdeps.Items, func(i, j int) bool {
		return sdeps.Items[i].Name < sdeps.Items[j].Name
	})
	return toServingList(sdeps.Items, opt), nil
}

func (s ServingController) Update(namespace string, sdepID string, sdep *seldonv1.SeldonDeployment) (*seldonv1.SeldonDeployment, error) {
	// 1. compose the update & return
	if err := Compose(sdep); err != nil {
		log.Errorf("Failed to compose the Seldon Deployment: %v", err)
		return nil, errors.RenderError(err)
	}

	// 2. execute the update & return
	result, err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Update(context.TODO(), sdep, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.RenderError(err)
	}
	return result, nil
}

func (s ServingController) Delete(namespace, sdepID string) error {
	err := s.seldonClient.MachinelearningV1().SeldonDeployments(namespace).Delete(context.TODO(), sdepID, metav1.DeleteOptions{})
	if err != nil {
		return errors.RenderError(err)
	}
	return nil
}
