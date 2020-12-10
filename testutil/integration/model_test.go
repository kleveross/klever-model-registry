/*Case Range
Model Extract：
	DO Not Extract（TensorRT，非模型（待定））；
	Do Extract (SavedModel,ONNX,GraphDef,NetDef,Keras H5, CaffeModel, TorchScript, MXNetParams, PMML)-create ModelJob automatically
Model Convert:
	Do Not Convert
	Do Convert (MXNetParams -> ONNX,Keras H5 -> SavedModel,CaffeModel -> NetDef)
API
	modeljob list, create, get, delete, get events
*/
package integration

import (
	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/ormb/pkg/oras"
	"github.com/kleveross/ormb/pkg/ormb"

	"path/filepath"
	"github.com/spf13/viper"
	"os"
	"github.com/sirupsen/logrus"

	"fmt"
	"net/http"
	"time"

	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

var _ = Describe("Model Registry", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1

	var modelExtract = []struct {
		model   string
		version string
		format  modeljobsv1alpha1.Format
	}{

		{"onnx", "v1", modeljobsv1alpha1.FormatONNX},
		{"savedmodel", "v1", modeljobsv1alpha1.FormatSavedModel},
		{"graphdef", "v1", modeljobsv1alpha1.FormatGraphDef},
		{"netdef", "v1", modeljobsv1alpha1.FormatNetDef},
		{"h5", "v1", modeljobsv1alpha1.FormatH5},
		{"caffe", "v1", modeljobsv1alpha1.FormatCaffeModel},
		{"mxnetparams", "v1", modeljobsv1alpha1.FormatMXNETParams},
		{"torchscript", "v1", modeljobsv1alpha1.FormatTorchScript},
		{"pmml", "v1", modeljobsv1alpha1.FormatPMML},
		//{"tensorrt","v1", "TensorRT"},
	}

	var modelConvert = []struct {
		model   string
		tomodel string
		version string
		format  modeljobsv1alpha1.Format
		toformat modeljobsv1alpha1.Format

	}{
		{"h5", "tensorflow","v1", modeljobsv1alpha1.FormatH5, modeljobsv1alpha1.FormatSavedModel},
		{"caffe", "netdef","v1", modeljobsv1alpha1.FormatCaffeModel, modeljobsv1alpha1.FormatNetDef},
		{"mxnetparams", "onnx", "v1", modeljobsv1alpha1.FormatMXNETParams, modeljobsv1alpha1.FormatONNX},
	}

	e := httpexpect.New(GinkgoT(), ModelRegistryHost)

	Context("Models", func() {
		project := "library"
		//check the model status
		It("Should list the model successfully", func() {
			//Upload the model by ormb
			var tool ormb.Interface
			plainHTTPOpt := true
			var insecureOpt bool
			project := "library"
			rootPath, err := filepath.Abs(viper.GetString("./models"))
			Expect(err).Should(BeNil())
			logrus.WithFields(logrus.Fields{
				"root-path": rootPath,
			}).Debugln("Create the ormb client with the given root path")

			tool, err = ormb.New(
				oras.ClientOptRootPath(rootPath),
				oras.ClientOptWriter(os.Stdout),
				oras.ClientOptPlainHTTP(plainHTTPOpt),
				oras.ClientOptInsecure(insecureOpt),
				)
			for _, tt := range modelExtract{
				dirtag := fmt.Sprintf("%s/%s/%s:%s", ModelRegistryHost[7:], project, tt.model, tt.version)
				err = tool.Save("./models/"+tt.model, dirtag)
				Expect(err).Should(BeNil())
				err = tool.Push(dirtag)
				Expect(err).Should(BeNil())

				//check the models
				artifact := e.GET("/api/v2.0/projects/{project_name}/repositories/{repository_name}/artifacts/{version}",
					project, tt.model, tt.version).Expect().Status(http.StatusOK).JSON().Object()
				// Validate that the artifact is a SavedModel.
				// It is blocked since goharbor/harbor-helm does not support 2.1 now.
				// artifact.Value("extra_attrs").Object().Value("format").Equal("SavedModel")
				artifact.Value("type").Equal("MODEL")
			}
		})
		//create modeljob for convert by modelconvert
		Context("Convert ModelJobs", func() {
			It("Should create the ModelJobs successfully", func() {
				for _, tt := range modelConvert {
					name := tt.model+"2"+tt.tomodel
					desiredtag := fmt.Sprintf("%s/%s/%s:%s", ModelRegistryHost[7:], project, tt.tomodel, "new")
					job := modeljobsv1alpha1.ModelJob{
						ObjectMeta: metav1.ObjectMeta{
							Name: name,
						},
						Spec: modeljobsv1alpha1.ModelJobSpec{
							// Use ModelRegistryHost/<project>/<model>:<version>.
							Model: fmt.Sprintf("%s/%s/%s:%s", ModelRegistryHost[7:], project, tt.model, tt.version),
							DesiredTag: &desiredtag,
							ModelJobSource: modeljobsv1alpha1.ModelJobSource{
								Conversion: &modeljobsv1alpha1.ConversionSource{
									MMdnn: &modeljobsv1alpha1.MMdnnSpec{
										modeljobsv1alpha1.ConversionBaseSpec{
											From: tt.format,
											To: tt.toformat,
										},
									},
								},
							},
						},
					}
					e.POST("/api/v1alpha1/namespaces/{namespace}/modeljobs", "default").WithJSON(job).Expect().Status(http.StatusCreated)
				}
			})
			//check the modeljob is created
			It("Should list the ModelJobs successfully", func() {
				for _, tt := range modelConvert {
					name := tt.model+"2"+tt.tomodel
					modelJobs := e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/",
						"default").Expect().Status(http.StatusOK).JSON().Object().Value("items").Array()
					found := false
					for _, modelJob := range modelJobs.Iter() {
						rawLabels := modelJob.Path("$.metadata.name").Raw()
						if rawLabels.(string) == name {
							found = true
							// Set the name to the CRD name in the kubernetes cluster.
							name = modelJob.Path("$.metadata.name").String().Raw()
							return
						}
					}
					Expect(found).To(BeTrue())
				}
			})
			//check modeljob's status
			It("Should get the ModelJob successfully", func() {
				for _, tt := range modelConvert {
					name := tt.model + "2" + tt.tomodel
					e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}",
						"default", name).Expect().Status(http.StatusOK)
				}
			})
			// get modeljob event
			It("Should get the ModelJob events successfully", func() {
				for _, tt := range modelConvert {
					name := tt.model+"2"+tt.tomodel
					e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}/events",
						"default", name).Expect().Status(http.StatusOK)
				}
			})
		})

		Context("Extract ModelJobs", func() {
			//create modeljob for extract by modelextract
			It("Should create the ModelJobs successfully", func() {
				for _, tt := range modelExtract {
					name := tt.model + "extract"
					job := modeljobsv1alpha1.ModelJob{
						ObjectMeta: metav1.ObjectMeta{
							Name: name,
						},
						Spec: modeljobsv1alpha1.ModelJobSpec{
							// Use ModelRegistryHost/<project>/<model>:<version>.
							Model: fmt.Sprintf("%s/%s/%s:%s", ModelRegistryHost[7:], project, tt.model, tt.version),
							ModelJobSource: modeljobsv1alpha1.ModelJobSource{
								Extraction: &modeljobsv1alpha1.ExtractionSource{
									Format: tt.format,
								},
							},
						},
					}
					e.POST("/api/v1alpha1/namespaces/{namespace}/modeljobs", "default").WithJSON(job).Expect().Status(http.StatusCreated)
				}
			})
			//check the modeljob is created
			It("Should list the ModelJobs successfully", func() {
				for _, tt := range modelExtract{
					name := tt.model + "extract"
					modelJobs := e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/",
						"default").Expect().Status(http.StatusOK).JSON().Object().Value("items").Array()
					found := false
					for _, modelJob := range modelJobs.Iter() {
						rawLabels := modelJob.Path("$.metadata.name").Raw()
						if rawLabels.(string) == name {
							found = true
							// Set the name to the CRD name in the kubernetes cluster.
							name = modelJob.Path("$.metadata.name").String().Raw()
							return
						}
					}
					Expect(found).To(BeTrue())
				}
			})

			It("Should get the ModelJob successfully", func() {
				for _, tt := range modelExtract{
					name := tt.model + "extract"
					e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}",
						"default", name).Expect().Status(http.StatusOK)}

			})

			It("Should get the ModelJob events successfully", func() {
				for _, tt := range modelExtract {
					name := tt.model + "extract"
					e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}/events",
						"default", name).Expect().Status(http.StatusOK)
				}
			})
		})
	})
})
