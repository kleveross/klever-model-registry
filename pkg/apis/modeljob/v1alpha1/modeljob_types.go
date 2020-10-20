/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// FrameworkEnvKey is the env key of framework
	FrameworkEnvKey = "FRAMEWORK"
	// FormatEnvKey is the env key of format
	FormatEnvKey = "FORMAT"
	// ScriptCommandEnvKey is the python command of script
	ScriptCommandEnvKey = "COMMAND"
	// SourceModelPathEnvKey is the path env key of ormb pull
	SourceModelPathEnvKey = "SOURCE_MODEL_PATH"
	// DestinationModelPathEnvKey is the the path env key of convert's destination
	DestinationModelPathEnvKey = "DESTINATION_MODEL_PATH"
	// SourceModelTagEnvKey is the source tag of model
	SourceModelTagEnvKey = "SOURCE_MODEL_TAG"
	// DestinationModelTagEnvKey is the destination tag of model
	DestinationModelTagEnvKey = "DESTINATION_MODEL_TAG"
	// ExtractorEnvKey is extractor env key
	ExtractorEnvKey = "EXTRACTOR"

	// SourceModelPath is path of ormb pull
	SourceModelPath = "/models/input"
	// DestinationModelPath if the path of convert's destination
	DestinationModelPath = "/models/output"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ModelJob is the Schema for the modeljobs API
type ModelJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelJobSpec   `json:"spec,omitempty"`
	Status ModelJobStatus `json:"status,omitempty"`
}

// ModelJobSpec defines the desired state of ModelJob
type ModelJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Model is model ref, eg: kleveross/resnet:v1.
	Model string `json:"model,omitempty"`

	// DesiredTag is the target tag of model convert.
	DesiredTag *string `json:"desiredTag,omitempty"`

	// ModelJobSource is model job source.
	ModelJobSource `json:",inline"`
}

// ModelJobSource defines the modeljob source information
type ModelJobSource struct {
	Extraction *ExtractionSource `json:"extraction,omitempty"`
	Conversion *ConversionSource `json:"conversion,omitempty"`
}

type ExtractionSource struct {
	Format Format `json:"format,omitempty"`
}

type ConversionSource struct {
	MMdnn *MMdnnSpec `json:"mmdnn,omitempty"`
}

type MMdnnSpec struct {
	ConversionBaseSpec `json:",inline"`
}

type ConversionBaseSpec struct {
	From Format `json:"from,omitempty"`
	To   Format `json:"to,omitempty"`
}

// Framework is model framework, eg: TensorFlow.
type Framework string

const (
	FrameworkTensorflow Framework = "TensorFlow"
	FrameworkPyTorch    Framework = "PyTorch"
	FrameworkCaffe      Framework = "Caffe"
	FrameworkCaffe2     Framework = "Caffe2"
	FrameworkMXNet      Framework = "MXNet"
	FrameworkKeras      Framework = "Keras"
	FrameworkOthers     Framework = "Others"
	FrameworkONNX       Framework = "ONNX"
	FrameworkTensorRT   Framework = "TensorRT"
	FrameworkPMML       Framework = "PMML"
)

// Format is model format, eg: SaveModel.
type Format string

const (
	FormatSavedModel  Format = "SavedModel"
	FormatONNX        Format = "ONNX"
	FormatH5          Format = "H5"
	FormatPMML        Format = "PMML"
	FormatCaffeModel  Format = "CaffeModel"
	FormatNetDef      Format = "NetDef"
	FormatMXNETParams Format = "MXNETParams"
	FormatTorchScript Format = "TorchScript"
	FormatGraphDef    Format = "GraphDef"
	FormatTensorRT    Format = "TensorRT"
	FormatSKLearn     Format = "SKLearn"
	FormatXGBoost     Format = "XGBoost"
	// TODO: we have exactly same data structure in ormb
	FormatMLflow Format = "MLflow"
)

type ModelJobPhase string

const (
	ModelJobEmpty     ModelJobPhase = ""
	ModelJobPending   ModelJobPhase = "Pending"
	ModelJobRunning   ModelJobPhase = "Running"
	ModelJobDeleting  ModelJobPhase = "Deleting"
	ModelJobSucceeded ModelJobPhase = "Succeeded"
	ModelJobFailed    ModelJobPhase = "Failed"
)

// ModelJobStatus defines the observed state of ModelJob
type ModelJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ModelJobPhase is model status.
	Phase ModelJobPhase `json:"phase"`

	// Human readable message indicating the reason for Failure
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ModelJobList contains a list of ModelJob
type ModelJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelJob{}, &ModelJobList{})
}
